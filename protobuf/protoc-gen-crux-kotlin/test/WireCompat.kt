// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

// Compares the Kotlin message classes with the Java lite classes of the same proto files. For
// each sample from gensamples, both classes parse the sample and write it again, and the bytes
// must be equal in both directions.
import java.io.File
import kotlin.system.exitProcess

private fun parse(cls: Class<*>, bytes: ByteArray): Any =
    cls.getMethod("parseFrom", ByteArray::class.java).invoke(null, bytes)

private fun bytesOf(message: Any): ByteArray = message.javaClass.getMethod("toByteArray").invoke(message) as ByteArray

private fun sizeOf(message: Any): Int = message.javaClass.getMethod("getSerializedSize").invoke(message) as Int

fun main(args: Array<String>) {
    val dir = File(args[0])
    var failures = 0
    var samples = 0
    var identical = 0
    var openEnumReordered = 0
    val covered = sortedSetOf<String>()
    for (line in File(dir, "index.tsv").readLines()) {
        val (kotlinName, javaName, file, fullName, openEnumText) = line.split('\t')
        val openEnum = openEnumText == "true"
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
            val javaAgain = parse(javaClass, kBytes)
            if (sizeOf(k) != kBytes.size) fail("Kotlin getSerializedSize ${sizeOf(k)} but toByteArray ${kBytes.size}")
            if (openEnum) {
                // Java lite keeps an unknown enum number in its field; the Kotlin enum cannot hold it, so
                // the number goes to the unknown fields and is written last.
                if (javaAgain != j) fail("Java lite parse of the Kotlin bytes is not equal to the parse of the sample")
                if (!kBytes.contentEquals(jBytes)) openEnumReordered++
            } else {
                if (!kBytes.contentEquals(jBytes)) fail("Kotlin writes ${kBytes.size} bytes, Java lite writes ${jBytes.size} bytes different")
                if (sizeOf(j) != sizeOf(k)) fail("getSerializedSize differs: Kotlin ${sizeOf(k)} Java ${sizeOf(j)}")
                if (!bytesOf(javaAgain).contentEquals(jBytes)) fail("Java lite re-parse of the Kotlin bytes differs")
            }
            val kAgain = parse(kotlinClass, jBytes)
            if (!openEnum && (kAgain != k || kAgain.hashCode() != k.hashCode())) fail("Kotlin parse of the Java bytes is not equal")
            if (!bytesOf(kAgain).contentEquals(kBytes)) fail("Kotlin re-write of the Java bytes differs")
            if (kBytes.contentEquals(input)) identical++
            covered += fullName
            samples++
        } catch (e: Throwable) {
            fail(e.toString())
        }
    }
    println("samples=$samples messages=${covered.size} identicalToInput=$identical openEnumReordered=$openEnumReordered failures=$failures")
    exitProcess(if (failures == 0 && samples > 0) 0 else 1)
}
