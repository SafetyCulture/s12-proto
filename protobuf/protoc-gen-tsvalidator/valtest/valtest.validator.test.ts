import {
  ValTestMessageValidator,
  ValTestMessage_Email,
  ValTestMessage_Name,
  ValTestMessage_Description,
} from "./valtest.validator";
import fs from "fs";

describe("validate files", () => {
  test("validate names", () => {
    const inputs = readFiles(["file:///testdata/valid_names.txt"]);
    inputs.forEach((input) => {
      expect(() => ValTestMessage_Name.parse(input)).not.toThrow();
    });
  });

  test("validate safe strings", () => {
    const inputs = readFiles(["file:///testdata/valid_safe_strings.txt"]);
    inputs.forEach((input) => {
      expect(() => ValTestMessage_Description.parse(input)).not.toThrow();
    });
  });

  test("validate invalid emails", () => {
    const inputs = readFiles([
      "file:///testdata/invalid_emails.txt",
      // "example|2@example.com", // TODO: Zod doesnt like
      // "name@組織.香港" // TODO: Zod doesnt like
    ]);
    inputs.forEach((input) => {
      expect(() => ValTestMessage_Email.parse(input), input).toThrow();
    });
  });

  test("validate valid emails", () => {
    const inputs = readFiles([
      "file:///testdata/valid_emails.txt",
      // "valid.aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa64@aa256aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa.com" // TODO: Zod doesnt like
    ]);
    inputs.forEach((input) => {
      expect(() => ValTestMessage_Email.parse(input), input).not.toThrow();
    });
  });

  // TODO: Fix this test
  // test("validate invalid safe strings", () => {
  //   const inputs = readFiles(["file:///invalid_safe_strings.txt"]);
  //   inputs.forEach((input) => {
  //     expect(() => ValTestMessage_Description.parse(input), input).toThrow();
  //   });
  // });
});

describe("validate message", () => {
  test("ValTestMessage - valMsgOpts", () => {
    expect(() => ValTestMessageValidator.parse(valMsgOpts)).not.toThrow();
  });

  test("ValTestMessage - valMsg", () => {
    expect(() => ValTestMessageValidator.parse(getValMsg({}))).not.toThrow();
  });
});

const id = "92b6c2f9-abd8-48bc-a2c9-bf70e969751a";
const legacyId = "56341C6E-35A7-4C97-9C5E-7AC79673EAB2";
const s12Id = "audit_f6dad1c9334040739b1e67ca70f4cf4d";
// const legacyLongIdFail =
//   "00EAE67E-2160-4C2E-BEB1-E5558A2696A7-9-00000190327E0675"; // length = 49 (without dashes)
// const legacyLongId1 =
//   "00EAE67E-2160-4C2E-BEB1-E5558A2696A7-90-00000190327E0675"; // length = 50 (without dashes)
// const legacyLongId2 =
//   "005F2E38-8426-48AF-94DE-5FEA3A396EEA-891-00000153F68896DC"; // length = 51 (without dashes)
// const legacyLongId3 =
//   "007B516E-53F1-4AA0-ABAF-8C78342A2C82-2388-00000221F1C2BD1E"; // length = 52 (without dashes)
// const legacyLongId4 =
//   "00709A17-151F-4CFC-B412-F080343ED84D-11977-000010227B4C60A9"; // length = 53 (without dashes)
const email = "email@example.com";
const password = "1234567!";
// const name = "Γιώργος Νταλάρας";
// const idv4 = "92b6c2f9-abd8-48bc-a2c9-bf70e969751a";
// const uuidNotV4 = "e07b6ac0-8a05-11e2-9951-ddd1182f65d9";
// const valid = false;
// const invalid = true;
// const emptyString = "<EMPTY>";

const valMsg = {
  id: id,
  ids: [id, id],
  legacyId: legacyId,
  mediaId: id,
  innerLegacyId: { id },
  innerS12Id: { id: s12Id },
  uuid: id,
  email: email,
  optEmail: email,
  description: "François Truffaut 久保田 利伸 text",
  password: password,
  title: "Short text, ok",
  fixedString: "abcd",
  runeString: "çois",
  replaceString: "With 'unsafe' \"<chars>",
  notReplaceString: "Safe chars ç only 123.",
  allowString: "Accept ~ and #",
  symbolString: "Accept $ £ ¥ €",
  symbolsString: "Accept 🌏 💯 and a\u030C", // ǎ as in a + the caron
  newlineString: "Accept\nNewlines\n\rYeeha",
  invalidEncodingString: "Accept invalid \xe9",
  optString: "Optional",
  trimString: "   Trim me   \t",
  allString: " Lot of checks here>",
  name: "Sinéad O'Connor",
  noValidation: "<really?>' OR 1=1",
  contactOneof: { phone: "14574560123" },
  msgRequired: { id: id },
  nestedMessage: {
    val: "inner val",
    // innerNestedMessage: {
    // 	innerVal: "abc def",
    // },
    nestedEmail: email,
    memberEmails: [email, email],
  },
  contactsWithLengthConstraint: [
    { phone: "abc", email: "test@example.com" },
    { phone: "", email: "test2@example.com" },
  ],
  url: "https://example.com/test",
  urlAllOpts: "http://app.safetyculture.com/report/media?param=test#fragment",
  s12Id: id, // Diff = s12Id
  phone: "14574560123", // Diff
  notSupported: {}, // Diff
  timezone: "Australia/Sydney",
  longString: "x".repeat(30000),
};

const getValMsg = (override: any) => {
  return { ...valMsg, ...override };
};

// omit optional fields here
const valMsgOpts = {
  id: id,
  ids: [id, id],
  legacyId: legacyId,
  innerLegacyId: { id: legacyId },
  innerS12Id: { id: s12Id },
  uuid: id,
  email: email,
  description: "François Truffaut 久保田 利伸 text",
  password: password,
  title: "Short text, ok",
  fixedString: "abcd",
  runeString: "çois",
  replaceString: "With 'unsafe' \"<chars>",
  notReplaceString: "Safe chars ç only 123.",
  allowString: "Accept ~ and #",
  symbolString: "Accept $ £ ¥ €",
  symbolsString: "Accept 🌏 💯 and a\u030C", // ǎ as in a + the caron
  newlineString: "Accept\nNewlines\n\rYeeha",
  invalidEncodingString: "Accept invalid \xe9",
  trimString: "   Trim me   \t",
  allString: " Lot of checks here>",
  noValidation: "<really?>' OR 1=1",
  contactOneof: { phone: "14574560123" },
  msgRequired: { id: id },
  nestedMessage: {
    val: "inner val",
    // innerNestedMessage: {
    // 	innerVal: "abc def",
    // },
    nestedEmail: email,
    memberEmails: [email, email],
  },
  contactsWithLengthConstraint: [{ phone: "abc", email: "test@example.com" }],
  url: "https://example.com/test",
  s12Id: id, // Diff = s12Id

  phone: "14574560123", // Diff
  notSupported: {}, // Diff
  timezone: "Australia/Sydney",
};

function readFiles(list: string[]): string[] {
  const outList: string[] = [];
  for (const filename of list) {
    if (!filename.startsWith("file:///")) {
      outList.push(filename);
      continue;
    }
    const path = process.cwd();
    const filePath = filename.replace("file:///", `${path}/`);
    try {
      const fileContents = fs.readFileSync(filePath, "utf-8");
      const lines = fileContents.split("\n");
      outList.push(...lines);
    } catch (err) {
      console.error(`Cannot open: ${err}`);
    }
  }
  return outList;
}
