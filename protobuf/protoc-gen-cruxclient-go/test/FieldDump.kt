// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

// Writes every field of a Kotlin message and of a Java lite message in one text form, by reflection
// and with the manifest from gensamples, so that the test can compare the decoded values.
import java.io.File

class FieldInfo(
    val number: Int,
    val property: String,
    val javaName: String,
    val shape: String,
    val kind: String,
    val openEnum: Boolean,
    val messageClass: String,
    val oneofProperty: String,
    val oneofJavaName: String,
    val oneofMember: String,
)

fun readManifest(file: File): Map<String, List<FieldInfo>> {
    val result = linkedMapOf<String, MutableList<FieldInfo>>()
    var current: MutableList<FieldInfo>? = null
    for (line in file.readLines()) {
        val c = line.split('\t')
        if (c[0] == "M") {
            current = mutableListOf()
            result[c[1]] = current
        } else {
            checkNotNull(current) { "field line before the first M line: $line" }.add(FieldInfo(c[1].toInt(), c[2], c[3], c[4], c[5], c[6] == "1", c[7], c[8], c[9], c[10]))
        }
    }
    return result
}

private fun call(target: Any, name: String): Any? =
    try {
        target.javaClass.getMethod(name).invoke(target)
    } catch (e: java.lang.reflect.InvocationTargetException) {
        throw e.targetException
    }

/** Calls an accessor of a Java lite message. A field whose name clashes with a method of a base class gets an underscore. */
private fun javaAccessor(target: Any, prefix: String, name: String, suffix: String = ""): Any? {
    val plain = prefix + name + suffix
    val own = target.javaClass.declaredMethods.any { it.name == plain && it.parameterCount == 0 }
    return if (own) call(target, plain) else call(target, prefix + name + "_" + suffix)
}

private fun kotlinGetter(target: Any, property: String): Any? {
    val getter = "get" + property.replaceFirstChar { it.uppercase() }
    return if (target.javaClass.methods.any { it.name == getter && it.parameterCount == 0 }) call(target, getter) else call(target, property)
}

private fun hex(bytes: ByteArray) = bytes.joinToString("") { (it.toInt() and 0xFF).toString(16).padStart(2, '0') }

private fun scalar(kind: String, value: Any?): String =
    when (kind) {
        "float" -> (value as Float).toRawBits().toString()
        "double" -> (value as Double).toRawBits().toString()
        else -> value.toString()
    }

private fun present(f: FieldInfo, value: Any?): Any = checkNotNull(value) { "field #${f.number} has a null ${f.kind} value" }

private fun kotlinEnumValue(value: Any): String {
    val name = call(value, "getName").toString()
    return if (name == "UNRECOGNIZED") call(value, "getNumber").toString() else name
}

private fun javaEnumValue(value: Any, number: Int? = null): String {
    val name = call(value, "name").toString()
    return if (name == "UNRECOGNIZED") number.toString() else name
}

private fun kotlinValue(manifest: Map<String, List<FieldInfo>>, f: FieldInfo, value: Any?): String =
    when (f.kind) {
        "bytes" -> hex(call(present(f, value), "toByteArray") as ByteArray)
        "enum" -> kotlinEnumValue(present(f, value))
        "message" -> dumpKotlin(manifest, present(f, value))
        else -> scalar(f.kind, value)
    }

private fun javaValue(manifest: Map<String, List<FieldInfo>>, f: FieldInfo, value: Any?, number: Int? = null): String =
    when (f.kind) {
        "bytes" -> hex(call(present(f, value), "toByteArray") as ByteArray)
        "enum" -> javaEnumValue(present(f, value), number)
        "message" -> dumpJava(manifest, present(f, value), f.messageClass)
        else -> scalar(f.kind, value)
    }

fun dumpKotlin(manifest: Map<String, List<FieldInfo>>, message: Any): String {
    val fields = manifest.getValue(message.javaClass.name)
    val out = StringBuilder("{")
    val oneofDone = mutableSetOf<String>()
    for (f in fields) {
        val value = kotlinGetter(message, if (f.shape == "oneof") f.oneofProperty else f.property)
        out.append(' ')
        when (f.shape) {
            "single" -> out.append("#${f.number}=").append(kotlinValue(manifest, f, value))
            "optional" -> out.append("#${f.number}=").append(if (value == null) "<unset>" else kotlinValue(manifest, f, value))
            "list" -> out.append("#${f.number}=[").append((value as List<*>).joinToString(",") { kotlinValue(manifest, f, it) }).append("]")
            "map" -> out.append("#${f.number}={").append((value as Map<*, *>).entries.map { "${it.key}:${kotlinValue(manifest, f, it.value)}" }.sorted().joinToString(",")).append("}")
            else -> {
                if (!oneofDone.add(f.oneofProperty)) {
                    out.setLength(out.length - 1)
                    continue
                }
                if (value == null) {
                    out.append("oneof ${f.oneofProperty}=<none>")
                } else {
                    val member = fields.first { it.oneofProperty == f.oneofProperty && it.oneofMember == value.javaClass.simpleName }
                    out.append("oneof ${f.oneofProperty}=${member.number}:").append(kotlinValue(manifest, member, call(value, "getValue")))
                }
            }
        }
    }
    return out.append(" }").toString()
}

fun dumpJava(manifest: Map<String, List<FieldInfo>>, message: Any, kotlinClass: String): String {
    val fields = manifest.getValue(kotlinClass)
    val out = StringBuilder("{")
    val oneofDone = mutableSetOf<String>()
    for (f in fields) {
        val suffix = ""
        out.append(' ')
        when (f.shape) {
            "single" -> out.append("#${f.number}=").append(javaValue(manifest, f, javaAccessor(message, "get", f.javaName, suffix), if (f.kind == "enum" && f.openEnum) javaAccessor(message, "get", f.javaName, "Value") as Int else null))
            "optional" ->
                out.append("#${f.number}=").append(
                    if (javaAccessor(message, "has", f.javaName) == false) "<unset>" else javaValue(manifest, f, javaAccessor(message, "get", f.javaName, suffix), if (f.kind == "enum" && f.openEnum) javaAccessor(message, "get", f.javaName, "Value") as Int else null),
                )
            "list" -> {
                val values = javaAccessor(message, "get", f.javaName, "${suffix}List") as List<*>
                val numbers = if (f.kind == "enum" && f.openEnum) javaAccessor(message, "get", f.javaName, "ValueList") as List<Int> else emptyList()
                out.append("#${f.number}=[").append(values.mapIndexed { index, value -> javaValue(manifest, f, value, numbers.getOrNull(index)) }.joinToString(",")).append("]")
            }
            "map" -> {
                val values = javaAccessor(message, "get", f.javaName, "${suffix}Map") as Map<*, *>
                val numbers = if (f.kind == "enum" && f.openEnum) javaAccessor(message, "get", f.javaName, "ValueMap") as Map<*, *> else emptyMap<Any?, Any?>()
                out.append("#${f.number}={").append(values.entries.map { "${it.key}:${javaValue(manifest, f, it.value, numbers[it.key] as? Int)}" }.sorted().joinToString(",")).append("}")
            }
            else -> {
                if (!oneofDone.add(f.oneofProperty)) {
                    out.setLength(out.length - 1)
                    continue
                }
                val case = checkNotNull(call(message, "get${f.oneofJavaName}Case")) { "null case of oneof ${f.oneofJavaName}" }
                val number = call(case, "getNumber") as Int
                if (number == 0) {
                    out.append("oneof ${f.oneofProperty}=<none>")
                } else {
                    val member = fields.first { it.oneofProperty == f.oneofProperty && it.number == number }
                    val memberSuffix = ""
                    val memberNumber = if (member.kind == "enum" && member.openEnum) javaAccessor(message, "get", member.javaName, "Value") as Int else null
                    out.append("oneof ${f.oneofProperty}=$number:").append(javaValue(manifest, member, javaAccessor(message, "get", member.javaName, memberSuffix), memberNumber))
                }
            }
        }
    }
    return out.append(" }").toString()
}

/** Returns the text around the first position where two dumps differ. */
fun firstDifference(a: String, b: String): String {
    val i = a.commonPrefixWith(b).length
    val from = maxOf(0, i - 40)
    return "Kotlin ...${a.substring(from, minOf(a.length, i + 60))} / Java ...${b.substring(from, minOf(b.length, i + 60))}"
}
