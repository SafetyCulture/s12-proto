// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

namespace S12.Protobuf.Validation;

/// <summary>Receives violations of rules marked <c>log_only</c>, which never fail validation.</summary>
public static class ValidationLog
{
    /// <summary>
    /// Invoked with the field path, the unmet requirement, and the base64 of the field's first
    /// 50 characters. Discards the violation while unset.
    /// </summary>
    public static Action<string, string, string>? Handler { get; set; }

    public static void Report(string field, string requirement, string base64Value) =>
        Handler?.Invoke(field, requirement, base64Value);
}
