import { ValTestMessage_NestedMessage_NestedEmail } from "./valtest.validator";

describe("validate email", () => {
  test("validate email", () => {
    expect(() =>
      ValTestMessage_NestedMessage_NestedEmail.parse("a@msn.com"),
    ).not.toThrow();
    expect(() =>
      ValTestMessage_NestedMessage_NestedEmail.parse("amsn.com"),
    ).toThrow();
  });
});
