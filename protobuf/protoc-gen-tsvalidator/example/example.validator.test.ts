import {ExampleMessage_Contact_Email} from "./example.validator"

describe('validate email', () => {
  test('validate email', () => {
    expect(()=>ExampleMessage_Contact_Email.parse("a@msn.com")).not.toThrow()
    expect(()=>ExampleMessage_Contact_Email.parse("amsn.com")).toThrow()
  });
});