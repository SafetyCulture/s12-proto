// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

// Compares the Kotlin message classes with the Java lite classes of the same proto files. For
// each sample from gensamples, both classes parse the sample, every field must have the same
// value, and both classes write it again to equal bytes in both directions. Two consecutive samples
// of one message are also joined, in both orders: the join is a merge in which a singular field
// takes the last value, a oneof takes the last member, a repeated field appends and a message field
// merges. The kitchen checks at the end cover edits, nesting limits, merges, encodings and required
// fields.
import com.safetyculture.s12kt.kitchen.editions.ClosedMaps
import com.safetyculture.s12kt.kitchen.editions.Utf8Modes
import com.safetyculture.s12kt.kitchen.legacy.Legacy
import com.safetyculture.s12kt.kitchen.legacy.LevelMaps
import com.safetyculture.s12kt.kitchen.legacy.Strict
import com.safetyculture.s12kt.kitchen.legacy.Req
import com.safetyculture.s12kt.kitchen.legacy.ReqEnum
import com.safetyculture.s12kt.kitchen.legacy.ReqHolder
import com.safetyculture.s12kt.kitchen.legacy.RequiredMessage
import com.safetyculture.s12kt.kitchen.v1.Color
import com.safetyculture.s12kt.kitchen.v1.Maps
import com.safetyculture.s12kt.kitchen.v1.OptionalScalars
import com.safetyculture.s12kt.kitchen.v1.Oneofs
import com.safetyculture.s12kt.kitchen.v1.Repeated
import com.safetyculture.s12kt.kitchen.v1.Scalars
import com.safetyculture.s12kt.wire.InvalidWireException
import com.safetyculture.s12kt.wire.UninitializedMessageException
import com.safetyculture.s12kt.wire.WireReader
import java.io.ByteArrayOutputStream
import java.io.File
import java.lang.reflect.InvocationTargetException
import kotlin.system.exitProcess

private fun parse(cls: Class<*>, bytes: ByteArray): Any =
    try {
        cls.getMethod("parseFrom", ByteArray::class.java).invoke(null, bytes)
    } catch (e: InvocationTargetException) {
        throw e.targetException
    }

private fun bytesOf(message: Any): ByteArray = message.javaClass.getMethod("toByteArray").invoke(message) as ByteArray

private fun sizeOf(message: Any): Int = message.javaClass.getMethod("getSerializedSize").invoke(message) as Int

private fun compareSamples(dir: File): Int {
    var failures = 0
    var samples = 0
    var identical = 0
    var rejected = 0
    var merges = 0
    val covered = sortedSetOf<String>()
    val expected = File(dir, "sample-count.txt").readText().trim().toInt()
    var records = 0
    var previous: Pair<String, ByteArray>? = null
    for (line in File(dir, "index.tsv").readLines()) {
        records++
        val (kotlinName, javaName, file, fullName) = line.split('\t')
        val input = File(dir, file).readBytes()
        fun fail(what: String) {
            failures++
            println("FAIL $fullName $file: $what")
        }
        try {
            val kotlinClass = Class.forName(kotlinName)
            val javaClass = Class.forName(javaName)
            val kAccepted = accepts({ parse(kotlinClass, it) }, input)
            val jAccepted = accepts({ parse(javaClass, it) }, input)
            if (kAccepted == null || jAccepted == null || kAccepted != jAccepted) {
                fail("Kotlin accepted=$kAccepted, Java lite accepted=$jAccepted")
                continue
            }
            if (!kAccepted) {
                rejected++
                continue
            }
            val k = parse(kotlinClass, input)
            val j = parse(javaClass, input)
            val kBytes = bytesOf(k)
            val jBytes = bytesOf(j)
            val kDump = dumpKotlin(manifest, k)
            val jDump = dumpJava(manifest, j, kotlinName)
            if (kDump != jDump) fail("decoded values differ: ${firstDifference(kDump, jDump)}")
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
            val before = previous
            previous = Pair(kotlinName, input)
            if (before != null && before.first == kotlinName) {
                merges += mergeChecks(kotlinClass, javaClass, before.second, input, "$fullName $file", ::fail)
            }
        } catch (e: Throwable) {
            fail(e.toString())
        }
    }
    if (records != expected) {
        failures++
        println("FAIL sample index has $records records, gensamples wrote $expected")
    }
    println("samples=$records expected=$expected accepted=$samples messages=${covered.size} identicalToInput=$identical rejectedByBoth=$rejected merges=$merges failures=$failures")
    return failures + if (samples == 0) 1 else 0
}

/**
 * Joins two encoded messages in both orders and compares the parse of the join in Kotlin and in
 * Java lite: the same acceptance, the same value in every field and the same bytes. Returns the
 * number of joins that it compared.
 */
private fun mergeChecks(kotlin: Class<*>, java: Class<*>, a: ByteArray, b: ByteArray, label: String, fail: (String) -> Unit): Int {
    var compared = 0
    for ((first, second, order) in listOf(Triple(a, b, "a+b"), Triple(b, a, "b+a"))) {
        val joined = first + second
        val k = accepts({ parse(kotlin, it) }, joined)
        val j = accepts({ parse(java, it) }, joined)
        if (k == null || j == null || k != j) {
            fail("join $order: Kotlin accepted=$k, Java lite accepted=$j")
            continue
        }
        if (!k) continue
        val kMessage = parse(kotlin, joined)
        val jMessage = parse(java, joined)
        val kDump = dumpKotlin(manifest, kMessage)
        val jDump = dumpJava(manifest, jMessage, kotlin.name)
        if (kDump != jDump) fail("join $order of $label: values differ: ${firstDifference(kDump, jDump)}")
        if (!bytesOf(kMessage).contentEquals(bytesOf(jMessage))) fail("join $order of $label: written bytes differ")
        compared++
    }
    return compared
}

private typealias JScalars = com.safetyculture.s12.kitchen.v1.Scalars
private typealias JMaps = com.safetyculture.s12.kitchen.v1.Maps
private typealias JOptionalScalars = com.safetyculture.s12.kitchen.v1.OptionalScalars
private typealias JOneofs = com.safetyculture.s12.kitchen.v1.Oneofs
private typealias JRepeated = com.safetyculture.s12.kitchen.v1.Repeated
private typealias JLegacy = com.safetyculture.kitchen.legacy.Legacy
private typealias JReq = com.safetyculture.kitchen.legacy.Req
private typealias JReqEnum = com.safetyculture.kitchen.legacy.ReqEnum
private typealias JReqHolder = com.safetyculture.kitchen.legacy.ReqHolder
private typealias JLevelMaps = com.safetyculture.kitchen.legacy.LevelMaps
private typealias JStrict = com.safetyculture.kitchen.legacy.Strict
private typealias JClosedMaps = kitchen.editions.ClosedMaps
private typealias JUtf8Modes = kitchen.editions.Utf8Modes

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
        val start = System.nanoTime()
        val kotlin = accepts({ Oneofs.parseFrom(it) }, bytes)
        val millis = (System.nanoTime() - start) / 1_000_000
        val java = accepts({ JOneofs.parseFrom(it) }, bytes)
        failures += check("messages nested $depth deep: Kotlin=$kotlin Java=$java", kotlin != null && kotlin == java)
        // Each nested oneof message is parsed once, so the parse time grows with the depth, not with 2^depth.
        failures += check("messages nested $depth deep parse in under 1 s: ${millis} ms", millis < 1000)
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

private lateinit var manifest: Map<String, List<FieldInfo>>

/** Like [same], and an accepted message must also have the same value in every field. */
private fun sameValues(name: String, kotlin: Class<*>, java: Class<*>, bytes: ByteArray, accepted: Boolean): Int {
    var failures = same(name, { parse(kotlin, it) }, { parse(java, it) }, bytes, accepted)
    if (failures == 0 && accepted) {
        val kDump = dumpKotlin(manifest, parse(kotlin, bytes))
        val jDump = dumpJava(manifest, parse(java, bytes), kotlin.name)
        failures += check("$name: values", kDump == jDump)
        if (kDump != jDump) println("  ${firstDifference(kDump, jDump)}")
    }
    return failures
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
    holder("a map value that is missing, with a required field", delimited(3, delimited(1, "k".toByteArray())), false)
    failures += same("ReqEnum with a declared number", { ReqEnum.parseFrom(it) }, { JReqEnum.parseFrom(it) }, field(1, 0, varint(2)), true)
    failures += same("ReqEnum with an undeclared number", { ReqEnum.parseFrom(it) }, { JReqEnum.parseFrom(it) }, field(1, 0, varint(99)), false)
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

private fun fixed64(number: Int, bits: Long): ByteArray =
    field(number, 1, ByteArray(8) { (bits ushr (8 * it)).toByte() })

private fun fixed32(number: Int, bits: Int): ByteArray =
    field(number, 5, ByteArray(4) { (bits ushr (8 * it)).toByte() })

private fun vint(number: Int, value: Long): ByteArray {
    val out = ByteArrayOutputStream()
    var v = value
    while (v and 0x7FL.inv() != 0L) {
        out.write(((v and 0x7F) or 0x80).toInt())
        v = v ushr 7
    }
    out.write(v.toInt())
    return varint((number shl 3) or 0) + out.toByteArray()
}

/** The checks use one case for each way in which a value can arrive that the random samples do not produce. */
private fun valueChecks(): Int {
    var failures = 0
    fun scalars(name: String, bytes: ByteArray) {
        failures += sameValues("Scalars $name", Scalars::class.java, JScalars::class.java, bytes, true)
    }
    fun repeated(name: String, bytes: ByteArray) {
        failures += sameValues("Repeated $name", Repeated::class.java, JRepeated::class.java, bytes, true)
    }
    fun oneofs(name: String, bytes: ByteArray) {
        failures += sameValues("Oneofs $name", Oneofs::class.java, JOneofs::class.java, bytes, true)
    }
    fun optional(name: String, bytes: ByteArray) {
        failures += sameValues("OptionalScalars $name", OptionalScalars::class.java, JOptionalScalars::class.java, bytes, true)
    }

    val nanPayload = 0x7FF8000000000001L
    scalars("double NaN with a payload", fixed64(1, nanPayload))
    scalars("double negative zero", fixed64(1, Long.MIN_VALUE))
    scalars("double infinity", fixed64(1, 0x7FF0000000000000L))
    scalars("float NaN with a payload", fixed32(2, 0x7FC00001))
    scalars("float negative zero", fixed32(2, Int.MIN_VALUE))
    scalars("both zero", fixed64(1, 0L) + fixed32(2, 0))
    optional("optional double NaN", fixed64(1, nanPayload))
    optional("optional float negative zero", fixed32(2, Int.MIN_VALUE))
    repeated("packed doubles with NaN and zero", delimited(1, ByteArray(0) + ByteArray(8) { (nanPayload ushr (8 * it)).toByte() } + ByteArray(8)))
    repeated("unpacked doubles with NaN", fixed64(20, nanPayload) + fixed64(20, 0L))

    scalars("a singular field twice, the last value wins", vint(3, 5) + vint(3, 7))
    scalars("a singular field twice, the last value is the default", vint(3, 5) + vint(3, 0))
    scalars("a string twice", delimited(14, "a".toByteArray()) + delimited(14, "b".toByteArray()))
    scalars("a negative int32 in ten bytes", vint(3, -1L))
    scalars("an int32 with a value above 32 bits", vint(3, 0x1_0000_0005L))
    scalars("a bool that is 2", vint(13, 2))

    repeated("packed int32 in the unpacked form", vint(3, 1) + vint(3, 2) + vint(3, 300))
    repeated("unpacked int32 in the packed form", delimited(18, byteArrayOf(1, 2) + varint(300)))
    repeated("packed doubles in the unpacked form", fixed64(1, 1L) + fixed64(1, 2L))
    repeated("unpacked doubles in the packed form", delimited(20, ByteArray(8) { if (it == 0) 3 else 0 } + ByteArray(8)))
    repeated("packed enums in the unpacked form", vint(16, 1) + vint(16, 99) + vint(16, 2))
    repeated("unpacked enums in the packed form", delimited(19, byteArrayOf(1, 99, 2)))
    repeated("both forms mixed", delimited(3, byteArrayOf(1, 2)) + vint(3, 3) + delimited(3, byteArrayOf(4)))
    repeated("an empty packed run", delimited(3, ByteArray(0)) + vint(3, 9))

    val scalarsPart1 = vint(3, 5) + delimited(14, "x".toByteArray())
    val scalarsPart2 = vint(4, 8) + vint(3, 0)
    oneofs("a oneof member twice", delimited(2, "a".toByteArray()) + delimited(2, "b".toByteArray()))
    oneofs("a oneof member replaced by another member", delimited(2, "a".toByteArray()) + vint(3, 4))
    oneofs("a oneof message merged with itself", delimited(4, scalarsPart1) + delimited(4, scalarsPart2))
    oneofs("a oneof message replaced and then set again", delimited(4, scalarsPart1) + vint(3, 1) + delimited(4, scalarsPart2))
    oneofs("a oneof message with unknown fields in both parts", delimited(4, scalarsPart1 + vint(100, 1)) + delimited(4, scalarsPart2 + vint(101, 2)))
    failures += same(
        "oneof malformed message member then another member",
        { Oneofs.parseFrom(it) },
        { JOneofs.parseFrom(it) },
        delimited(4, delimited(14, byteArrayOf(0xFF.toByte()))) + delimited(2, "x".toByteArray()),
        false,
    )
    oneofs("two oneofs set", delimited(2, "a".toByteArray()) + vint(20, 1))
    optional("a message field merged with itself", delimited(17, scalarsPart1) + delimited(17, scalarsPart2))
    optional("a message field set to an empty message after a value", delimited(17, scalarsPart1) + delimited(17, ByteArray(0)))
    optional("an optional field set to its default", vint(3, 0))
    failures += sameValues(
        "Maps a key twice, the last value wins",
        Maps::class.java,
        JMaps::class.java,
        delimited(1, delimited(1, "k".toByteArray()) + delimited(2, "1".toByteArray())) + delimited(1, delimited(1, "k".toByteArray()) + delimited(2, "2".toByteArray())),
        true,
    )
    failures += sameValues(
        "Maps an entry with no key and no value",
        Maps::class.java,
        JMaps::class.java,
        delimited(1, ByteArray(0)) + delimited(8, ByteArray(0)),
        true,
    )
    failures += sameValues(
        "Maps a value message split in two entries",
        Maps::class.java,
        JMaps::class.java,
        delimited(8, vint(1, 1) + delimited(2, scalarsPart1)) + delimited(8, vint(1, 1) + delimited(2, scalarsPart2)),
        true,
    )
    failures += sameValues(
        "Maps an entry with two value messages",
        Maps::class.java,
        JMaps::class.java,
        byteArrayOf(0x42, 0x0D, 0x0D, 1, 0, 0, 0, 0x12, 2, 0x18, 1, 0x12, 2, 0x20, 2),
        true,
    )
    failures += same(
        "Maps a replaced message value with invalid UTF-8",
        { Maps.parseFrom(it) },
        { JMaps.parseFrom(it) },
        delimited(
            8,
            fixed32(1, 1) +
                delimited(2, delimited(14, byteArrayOf(0xFF.toByte()))) +
                delimited(2, scalarsPart1),
        ),
        false,
    )
    return failures
}

/** The entry of a map<string, enum>. */
private fun stringEntry(key: String, value: Long): ByteArray = delimited(1, key.toByteArray()) + vint(2, value)

private fun intEntry(key: Int, value: Long): ByteArray = vint(1, key.toLong()) + vint(2, value)

/**
 * Java lite filters the map entries with an enum number that a closed enum does not declare after
 * the whole message is read: it writes the entries again in the canonical form, after all other
 * unknown fields, and removes them from the map.
 */
private fun mapEnumChecks(): Int {
    var failures = 0
    fun closed(name: String, bytes: ByteArray) {
        failures += sameValues("ClosedMaps $name", ClosedMaps::class.java, JClosedMaps::class.java, bytes, true)
    }
    fun legacy(name: String, bytes: ByteArray) {
        failures += sameValues("LevelMaps $name", LevelMaps::class.java, JLevelMaps::class.java, bytes, true)
    }
    val unknownField = vint(9, 5)
    val overlongValue = byteArrayOf(0x99.toByte(), 0x80.toByte(), 0x00)
    for ((prefix, check) in listOf<Pair<String, (String, ByteArray) -> Unit>>("edition" to ::closed, "proto2" to ::legacy)) {
        check("$prefix: an entry with an unknown number", delimited(1, stringEntry("a", 99)))
        check("$prefix: an entry with a known number", delimited(1, stringEntry("a", 1)))
        check("$prefix: the same key, known and then unknown", delimited(1, stringEntry("k", 1)) + delimited(1, stringEntry("k", 99)))
        check("$prefix: the same key, unknown and then known", delimited(1, stringEntry("k", 99)) + delimited(1, stringEntry("k", 1)))
        check("$prefix: an unknown number and then an unknown field", delimited(1, stringEntry("a", 99)) + unknownField)
        check("$prefix: an unknown field and then an unknown number", unknownField + delimited(1, stringEntry("a", 99)))
        check("$prefix: two maps", delimited(2, intEntry(7, 99)) + delimited(1, stringEntry("b", 99)) + delimited(1, stringEntry("a", 99)))
        check("$prefix: a negative number", delimited(1, stringEntry("n", -1)))
        check("$prefix: a number above 32 bits", delimited(1, stringEntry("n", 0x1_0000_0005L)))
        check("$prefix: a value in more bytes than needed", delimited(1, delimited(1, "o".toByteArray()) + varint((2 shl 3) or 0) + overlongValue))
        check("$prefix: an entry without a value", delimited(1, delimited(1, "x".toByteArray())))
        check("$prefix: an entry with the value first", delimited(1, vint(2, 99) + delimited(1, "v".toByteArray())))
        check("$prefix: an entry with an unknown field inside", delimited(1, stringEntry("u", 99) + vint(5, 1)))
        check("$prefix: the same key three times", delimited(1, stringEntry("k", 99)) + delimited(1, stringEntry("k", 1)) + delimited(1, stringEntry("k", 98)))
    }
    closed("an open enum keeps the number in the map", delimited(3, stringEntry("a", 99)))
    closed("an entry with an unknown number and a field after it", delimited(1, stringEntry("a", 99)) + delimited(4, "tail".toByteArray()))
    return failures
}

/** Java lite checks the UTF-8 of a string for the Java options too: java_string_check_utf8 and the Java feature. */
private fun javaUtf8Checks(): Int {
    var failures = 0
    val invalid = byteArrayOf(0x61, 0xFF.toByte(), 0xC3.toByte(), 0x62)
    fun strict(name: String, bytes: ByteArray, accepted: Boolean) {
        failures += same("Strict $name", { Strict.parseFrom(it) }, { JStrict.parseFrom(it) }, bytes, accepted)
    }
    strict("optional string", delimited(1, invalid), false)
    strict("repeated string", delimited(2, "ok".toByteArray()) + delimited(2, invalid), false)
    strict("oneof string", delimited(3, invalid), false)
    strict("valid optional string", delimited(1, "\u00e9".toByteArray()), true)
    fun modes(name: String, bytes: ByteArray, accepted: Boolean) {
        failures += same("Utf8Modes $name", { Utf8Modes.parseFrom(it) }, { JUtf8Modes.parseFrom(it) }, bytes, accepted)
    }
    modes("string from the file feature", delimited(1, invalid), false)
    modes("string with the Java feature DEFAULT", delimited(2, invalid), true)
    modes("string with the core feature VERIFY", delimited(3, invalid), false)
    modes("string with the Java feature VERIFY", delimited(4, invalid), false)
    modes("nested string from the file feature", delimited(5, delimited(1, invalid)), false)
    return failures
}

/** Both classes must accept or both must refuse an input that is not a valid message. An accepted message must be equal. */
private fun malformedChecks(): Int {
    var failures = 0
    fun agree(name: String, bytes: ByteArray) {
        val k = accepts({ Scalars.parseFrom(it) }, bytes)
        val j = accepts({ JScalars.parseFrom(it) }, bytes)
        var ok = k != null && k == j
        if (ok && k == true) ok = bytesOf(Scalars.parseFrom(bytes)).contentEquals(bytesOf(JScalars.parseFrom(bytes)))
        failures += check("malformed $name: Kotlin=$k Java=$j", ok)
    }
    val ff = 0xFF.toByte()
    agree("truncated length", field(14, 2, varint(10) + "ab".toByteArray()))
    agree("length without data", varint((14 shl 3) or 2))
    agree("negative length", field(14, 2, byteArrayOf(ff, ff, ff, ff, 0x0F)))
    agree("huge length", field(14, 2, byteArrayOf(ff, ff, ff, 0x7F)))
    agree("ten byte varint", varint((3 shl 3) or 0) + ByteArray(9) { ff } + byteArrayOf(0x01))
    // Java lite reads an int32 or int64 varint of any length from a byte array; Kotlin refuses a varint of more than ten bytes.
    failures += check("malformed eleven byte varint is refused by Kotlin", accepts({ Scalars.parseFrom(it) }, varint((3 shl 3) or 0) + ByteArray(10) { ff } + byteArrayOf(0x01)) == false)
    failures += check("malformed eleven byte int64 varint is refused by Kotlin", accepts({ Scalars.parseFrom(it) }, varint((4 shl 3) or 0) + ByteArray(10) { ff } + byteArrayOf(0x01)) == false)
    agree("truncated varint", varint((3 shl 3) or 0) + byteArrayOf(ff))
    agree("truncated tag", byteArrayOf(0x80.toByte()))
    agree("a tag in ten bytes", ByteArray(9) { ff } + byteArrayOf(0x01))
    agree("a tag in six bytes", byteArrayOf(ff, ff, ff, ff, ff, 0x01))
    agree("wire type 6", varint((3 shl 3) or 6))
    agree("wire type 7", varint((3 shl 3) or 7))
    agree("wire type 6 and unknown number", varint((99 shl 3) or 6))
    agree("field number 0", byteArrayOf(0x00))
    agree("field number 0 with a value", field(0, 0, varint(1)))
    agree("a stray end group", varint((7 shl 3) or 4))
    agree("an end group of another number", varint((7 shl 3) or 3) + varint((8 shl 3) or 4))
    agree("a group that is not closed", varint((7 shl 3) or 3) + vint(1, 1))
    agree("a closed group", varint((7 shl 3) or 3) + vint(1, 1) + varint((7 shl 3) or 4))
    agree("a valid message and half a tag", vint(3, 5) + byteArrayOf(0x80.toByte()))
    agree("a valid message and a zero byte", vint(3, 5) + byteArrayOf(0x00))
    agree("fixed64 with seven bytes", varint((1 shl 3) or 1) + ByteArray(7))
    agree("fixed32 with three bytes", varint((2 shl 3) or 5) + ByteArray(3))
    agree("a declared field with the wrong wire type", varint((3 shl 3) or 2) + varint(1) + byteArrayOf(0x05))
    agree("a string with the varint wire type", vint(14, 5))
    agree("enum -1", vint(16, -1L))
    agree("enum -2147483648", vint(16, Int.MIN_VALUE.toLong()))
    agree("enum 4294967297", vint(16, 0x1_0000_0001L))
    agree("enum 4294967295", vint(16, 0xFFFF_FFFFL))
    for (size in listOf(16383, 16384, 16385, 70000, 1 shl 20)) {
        agree("string of $size bytes", delimited(14, ByteArray(size) { 'a'.code.toByte() }))
        agree("bytes of $size bytes", delimited(15, ByteArray(size) { it.toByte() }))
    }
    agree("a length that is one byte too long", delimited(14, ByteArray(300) { 'a'.code.toByte() }).let { it.copyOf(it.size - 1) })
    val big = ByteArray(40000) { (it % 7).toByte() }
    val repeatedBytes = delimited(18, big)
    failures += check(
        "malformed packed int32 of 40000 values",
        accepts({ Repeated.parseFrom(it) }, repeatedBytes) == accepts({ JRepeated.parseFrom(it) }, repeatedBytes) &&
            bytesOf(Repeated.parseFrom(repeatedBytes)).contentEquals(bytesOf(JRepeated.parseFrom(repeatedBytes))),
    )
    return failures
}

/** A caller can change a list after a size read without corrupting a later encode. */
private fun staleSizeChecks(): Int {
    val values = mutableListOf(1, 2, 3)
    val message = Repeated(fInt32 = values)
    message.getSerializedSize()
    values.addAll(listOf(300, 400, 500))
    return check(
        "a list changed after a size read writes the current list",
        JRepeated.parseFrom(message.toByteArray()).fInt32List == values,
    )
}

private fun requiredEncodeChecks(): Int {
    val missingMessage = try {
        RequiredMessage.parseFrom(ByteArray(0))
        false
    } catch (e: InvalidWireException) {
        true
    }
    val missingScalar = try {
        Req.decode(WireReader(reqNoName)).toByteArray()
        false
    } catch (e: UninitializedMessageException) {
        true
    }
    val explicitDefault = Req(name = "").toByteArray().contentEquals(JReq.newBuilder().setName("").build().toByteArray())
    var failures = check("a missing required message field is refused on decode", missingMessage)
    failures += check("a required scalar field that is missing is refused on encode", missingScalar)
    failures += check("a required scalar set to its default is encoded", explicitDefault)
    return failures
}

fun main(args: Array<String>) {
    manifest = readManifest(File(args[0], "fields.tsv"))
    var failures = compareSamples(File(args[0]))
    failures += editChecks()
    failures += nestingChecks()
    failures += requiredChecks()
    failures += utf8Checks()
    failures += valueChecks()
    failures += mapEnumChecks()
    failures += javaUtf8Checks()
    failures += requiredEncodeChecks()
    failures += malformedChecks()
    failures += staleSizeChecks()
    exitProcess(if (failures == 0) 0 else 1)
}
