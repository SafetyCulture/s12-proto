// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

// Compares the Kotlin message classes with the Java lite classes of the same proto files. For
// each sample from gensamples, both classes parse the sample and write it again, and the bytes
// must be equal in both directions. The kitchen checks at the end cover edits and nesting limits.
import com.safetyculture.s12kt.kitchen.legacy.Legacy
import com.safetyculture.s12kt.kitchen.legacy.Req
import com.safetyculture.s12kt.kitchen.legacy.ReqHolder
import com.safetyculture.s12kt.kitchen.v1.Color
import com.safetyculture.s12kt.kitchen.v1.Oneofs
import com.safetyculture.s12kt.kitchen.v1.Repeated
import com.safetyculture.s12kt.kitchen.v1.Scalars
import com.safetyculture.s12kt.wire.InvalidWireException
import java.io.ByteArrayOutputStream
import java.io.File
import kotlin.system.exitProcess

private fun parse(cls: Class<*>, bytes: ByteArray): Any =
    cls.getMethod("parseFrom", ByteArray::class.java).invoke(null, bytes)

private fun bytesOf(message: Any): ByteArray = message.javaClass.getMethod("toByteArray").invoke(message) as ByteArray

private fun sizeOf(message: Any): Int = message.javaClass.getMethod("getSerializedSize").invoke(message) as Int

private fun compareSamples(dir: File): Int {
    var failures = 0
    var samples = 0
    var identical = 0
    val covered = sortedSetOf<String>()
    for (line in File(dir, "index.tsv").readLines()) {
        val (kotlinName, javaName, file, fullName) = line.split('\t')
        val input = File(dir, file).readBytes()
        fun fail(what: String) {
            failures++
            println("FAIL $fullName $file: $what")
        }
        try {
            val kotlinClass = Class.forName(kotlinName)
            val javaClass = Class.forName(javaName)
            val k = parse(kotlinClass, input)
            val j = parse(javaClass, input)
            val kBytes = bytesOf(k)
            val jBytes = bytesOf(j)
            if (!kBytes.contentEquals(jBytes)) fail("Kotlin writes ${kBytes.size} bytes, Java lite writes ${jBytes.size} bytes different")
            if (sizeOf(k) != kBytes.size) fail("Kotlin getSerializedSize ${sizeOf(k)} but toByteArray ${kBytes.size}")
            if (sizeOf(j) != sizeOf(k)) fail("getSerializedSize differs: Kotlin ${sizeOf(k)} Java ${sizeOf(j)}")
            if (!bytesOf(parse(javaClass, kBytes)).contentEquals(jBytes)) fail("Java lite re-parse of the Kotlin bytes differs")
            val kAgain = parse(kotlinClass, jBytes)
            if (kAgain != k || kAgain.hashCode() != k.hashCode()) fail("Kotlin parse of the Java bytes is not equal")
            if (!bytesOf(kAgain).contentEquals(kBytes)) fail("Kotlin re-write of the Java bytes differs")
            if (kBytes.contentEquals(input)) identical++
            covered += fullName
            samples++
        } catch (e: Throwable) {
            fail(e.toString())
        }
    }
    println("samples=$samples messages=${covered.size} identicalToInput=$identical failures=$failures")
    return failures + if (samples == 0) 1 else 0
}

private typealias JScalars = com.safetyculture.s12.kitchen.v1.Scalars
private typealias JOneofs = com.safetyculture.s12.kitchen.v1.Oneofs
private typealias JRepeated = com.safetyculture.s12.kitchen.v1.Repeated
private typealias JLegacy = com.safetyculture.kitchen.legacy.Legacy
private typealias JReq = com.safetyculture.kitchen.legacy.Req
private typealias JReqHolder = com.safetyculture.kitchen.legacy.ReqHolder

private fun check(name: String, ok: Boolean): Int {
    println((if (ok) "ok   " else "FAIL ") + name)
    return if (ok) 0 else 1
}

/** An edit must replace an enum number that the enum does not declare, as it does in Java lite. */
private fun editChecks(): Int {
    var failures = 0

    val scalars = Scalars.parseFrom(JScalars.newBuilder().setFEnumValue(99).build().toByteArray())
    failures += check("a number that is not declared stays in place", scalars.fEnum == Color.UNRECOGNIZED(99) && scalars.fEnum.number == 99)
    failures += check("a number that is not declared is written again", JScalars.parseFrom(scalars.toByteArray()).fEnumValue == 99)
    val edited = scalars.copy(fEnum = Color.COLOR_RED)
    failures += check("copy replaces a single enum number", JScalars.parseFrom(edited.toByteArray()).fEnumValue == 1 && Scalars.parseFrom(edited.toByteArray()).fEnum == Color.COLOR_RED)

    val oneofs = Oneofs.parseFrom(JOneofs.newBuilder().setColorValue(99).build().toByteArray())
    failures += check("a oneof keeps a number that is not declared", oneofs.choice == Oneofs.Choice.Color(Color.UNRECOGNIZED(99)))
    val editedOneofs = oneofs.copy(choice = Oneofs.Choice.Name("x"))
    failures += check("copy replaces a oneof enum number", JOneofs.parseFrom(editedOneofs.toByteArray()).choiceCase == com.safetyculture.s12.kitchen.v1.Oneofs.ChoiceCase.NAME)

    val repeated = Repeated.parseFrom(JRepeated.newBuilder().addAllFEnumValue(listOf(1, 99, 2)).build().toByteArray())
    failures += check("a repeated enum keeps the order", JRepeated.parseFrom(repeated.toByteArray()).fEnumValueList == listOf(1, 99, 2))
    val editedRepeated = repeated.copy(fEnum = listOf(Color.COLOR_GREEN))
    failures += check("copy replaces a repeated enum list", JRepeated.parseFrom(editedRepeated.toByteArray()).fEnumValueList == listOf(2))
    return failures
}

private fun varint(value: Int): ByteArray {
    val out = ByteArrayOutputStream()
    var v = value
    while (v >= 0x80) {
        out.write((v and 0x7F) or 0x80)
        v = v ushr 7
    }
    out.write(v)
    return out.toByteArray()
}

/** Returns the bytes of an Oneofs message with [depth] more Oneofs messages inside it. */
private fun nestedMessages(depth: Int): ByteArray {
    var bytes = ByteArray(0)
    repeat(depth) { bytes = byteArrayOf(0xAA.toByte(), 0x01) + varint(bytes.size) + bytes }
    return bytes
}

/** Returns the bytes of [depth] nested groups of an unknown field. */
private fun nestedGroups(depth: Int): ByteArray {
    val start = varint((100 shl 3) or 3)
    val end = varint((100 shl 3) or 4)
    val out = ByteArrayOutputStream()
    repeat(depth) { out.write(start) }
    repeat(depth) { out.write(end) }
    return out.toByteArray()
}

private fun accepts(parser: (ByteArray) -> Any, bytes: ByteArray): Boolean? =
    try {
        parser(bytes)
        true
    } catch (e: InvalidWireException) {
        false
    } catch (e: com.google.protobuf.InvalidProtocolBufferException) {
        false
    } catch (e: Throwable) {
        println("unexpected ${e.javaClass.name}")
        null
    }

/** Java lite stops at a nesting depth of 100. Kotlin must stop at the same depth, with a checked error. */
private fun nestingChecks(): Int {
    var failures = 0
    for (depth in listOf(1, 50, 99, 100, 101, 102, 3000)) {
        val bytes = nestedMessages(depth)
        val kotlin = accepts({ Oneofs.parseFrom(it) }, bytes)
        val java = accepts({ JOneofs.parseFrom(it) }, bytes)
        failures += check("messages nested $depth deep: Kotlin=$kotlin Java=$java", kotlin != null && kotlin == java)
    }
    for (depth in listOf(1, 99, 100, 101, 200000)) {
        val bytes = nestedGroups(depth)
        val kotlin = accepts({ Scalars.parseFrom(it) }, bytes)
        val java = accepts({ JScalars.parseFrom(it) }, bytes)
        failures += check("unknown groups nested $depth deep: Kotlin=$kotlin Java=$java", kotlin != null && kotlin == java)
    }
    return failures
}

/** Parses the bytes with both classes. Both must accept or both must reject, and an accepted message must be written alike. */
private fun same(name: String, kotlin: (ByteArray) -> Any, java: (ByteArray) -> Any, bytes: ByteArray, accepted: Boolean): Int {
    val k = accepts(kotlin, bytes)
    val j = accepts(java, bytes)
    var ok = k == accepted && j == accepted
    if (ok && accepted) ok = bytesOf(kotlin(bytes)).contentEquals(bytesOf(java(bytes)))
    return check("$name: Kotlin=$k Java=$j", ok)
}

private fun field(number: Int, wireType: Int, payload: ByteArray = ByteArray(0)): ByteArray =
    varint((number shl 3) or wireType) + payload

private fun delimited(number: Int, payload: ByteArray): ByteArray = field(number, 2, varint(payload.size) + payload)

private val reqNoName = field(2, 0, varint(5))
private val reqWithName = delimited(1, "a".toByteArray())
private val legacyId = delimited(1, "id".toByteArray())
private val legacyScore = field(9, 1, ByteArray(8))

/** Java lite checks the required fields once, after it has merged the parts of a message. */
private fun requiredChecks(): Int {
    var failures = 0
    fun req(name: String, bytes: ByteArray, accepted: Boolean) {
        failures += same("Req $name", { Req.parseFrom(it) }, { JReq.parseFrom(it) }, bytes, accepted)
    }
    fun holder(name: String, bytes: ByteArray, accepted: Boolean) {
        failures += same("ReqHolder $name", { ReqHolder.parseFrom(it) }, { JReqHolder.parseFrom(it) }, bytes, accepted)
    }
    req("without the required field", reqNoName, false)
    req("with the required field", reqWithName + reqNoName, true)
    holder("a message field without the required field", delimited(1, reqNoName), false)
    holder("a message field split in two parts, the later part has the field", delimited(1, reqNoName) + delimited(1, reqWithName), true)
    holder("a message field split in two parts, the earlier part has the field", delimited(1, reqWithName) + delimited(1, reqNoName), true)
    holder("a message field split in two parts, no part has the field", delimited(1, reqNoName) + delimited(1, reqNoName), false)
    holder("a repeated element without the required field", delimited(2, reqWithName) + delimited(2, reqNoName), false)
    holder("a map value without the required field", delimited(3, delimited(1, "k".toByteArray()) + delimited(2, reqNoName)), false)
    holder("a map value with the required field", delimited(3, delimited(1, "k".toByteArray()) + delimited(2, reqWithName)), true)
    failures += same("Legacy without id", { Legacy.parseFrom(it) }, { JLegacy.parseFrom(it) }, legacyScore, false)
    failures += same("Legacy without score", { Legacy.parseFrom(it) }, { JLegacy.parseFrom(it) }, legacyId, false)
    failures += same("Legacy with both", { Legacy.parseFrom(it) }, { JLegacy.parseFrom(it) }, legacyId + legacyScore, true)
    failures += same("Legacy with a parent without id", { Legacy.parseFrom(it) }, { JLegacy.parseFrom(it) }, legacyId + legacyScore + delimited(8, legacyScore), false)
    return failures
}

/** Java lite checks the UTF-8 of a proto3 string and of a map entry. A proto2 string may hold invalid bytes. */
private fun utf8Checks(): Int {
    var failures = 0
    val invalid = byteArrayOf(0x61, 0xFF.toByte(), 0xC3.toByte(), 0x62)
    failures += same("proto2 string", { Legacy.parseFrom(it) }, { JLegacy.parseFrom(it) }, legacyId + legacyScore + delimited(3, invalid), true)
    failures += same("proto3 string", { Scalars.parseFrom(it) }, { JScalars.parseFrom(it) }, delimited(14, invalid), false)
    failures += same("proto3 string with valid bytes", { Scalars.parseFrom(it) }, { JScalars.parseFrom(it) }, delimited(14, "\u00e9".toByteArray()), true)
    failures += same("proto2 map key", { ReqHolder.parseFrom(it) }, { JReqHolder.parseFrom(it) }, delimited(4, delimited(1, invalid) + delimited(2, "x".toByteArray())), false)
    failures += same("proto2 map value", { ReqHolder.parseFrom(it) }, { JReqHolder.parseFrom(it) }, delimited(4, delimited(1, "k".toByteArray()) + delimited(2, invalid)), false)
    return failures
}

fun main(args: Array<String>) {
    var failures = compareSamples(File(args[0]))
    failures += editChecks()
    failures += nestingChecks()
    failures += requiredChecks()
    failures += utf8Checks()
    exitProcess(if (failures == 0) 0 else 1)
}
