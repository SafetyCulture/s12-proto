// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

namespace S12.Protobuf.Validation;

/// <summary>A message that can check its own fields against the rules declared in its .proto.</summary>
public interface IValidatableMessage
{
    /// <summary>
    /// Applies the declared rules, rewriting fields that carry mutating rules, and returns the first
    /// violation, or <c>null</c> when the message is valid.
    /// </summary>
    ValidationError? Validate();
}
