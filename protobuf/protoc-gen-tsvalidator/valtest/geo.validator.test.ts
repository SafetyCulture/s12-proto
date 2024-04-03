import {
  GeoValidationMessage_Latitude,
  GeoValidationMessage_Longitude,
  GeoValidationMessage_Accuracy,
  GeoValidationMessageValidator,
} from "./geo.validator";
import { z } from "zod";

const validLats = [0, 90, 90.0, 90.0, -0, 1.1, -1.1];
const validLatsWithUndefined = [...validLats, undefined];
const invalidLats = [-90.1, -91, 90.1, 91, null, "1"];

const validLons = [0, 180, 180.0, 180.0, -0, 1.1, -1.1];
const validLonsWithUndefined = [...validLons, undefined];
const invalidLons = [-180.1, -181, 180.1, 181, null, "1"];

const validAcc = [-0, 0, 1, 1.0, 10000];
const validAccWithUndefined = [...validAcc, undefined];
const invalidAcc = [-1, -1.0, 10001, null, "1"];

const validStrictMessages = [
  {},
  {
    latitude: 0,
  },
  {
    latitude: 0,
    longitude: 0,
  },
  {
    latitude: 0,
    longitude: 0,
    accuracy: 0,
  },
];
const invalidStrictMessages = [{ something: 1 }]
const validNonStrictMessages = [...validStrictMessages, ...invalidStrictMessages,]
const invalidMessages = [undefined];

describe("validate geo", () => {
  describe("Latitude", () => {
    test("valid number does not error", () => {
      validLatsWithUndefined.forEach((lat) => {
        expect(() => GeoValidationMessage_Latitude.parse(lat)).not.toThrow();
      });
    });

    test("valid string-number does not error on safeParse", () => {
      validLats.forEach((lat) => {
        expect(() =>
          z
            .preprocess((x) => Number(x), GeoValidationMessage_Latitude)
            .parse(JSON.stringify(lat)),
        ).not.toThrow();
      });
    });

    test("invalid errors", () => {
      invalidLats.forEach((lat) => {
        expect(() => GeoValidationMessage_Latitude.parse(lat)).toThrow();
      });
    });
  });

  describe("Latitude", () => {
    test("valid number does not error", () => {
      validLonsWithUndefined.forEach((lon) => {
        expect(() => GeoValidationMessage_Longitude.parse(lon)).not.toThrow();
      });
    });

    test("valid string-number does not error on safeParse", () => {
      validLons.forEach((lon) => {
        expect(() =>
          z
            .preprocess((x) => Number(x), GeoValidationMessage_Longitude)
            .parse(JSON.stringify(lon)),
        ).not.toThrow();
      });
    });

    test("invalid errors", () => {
      invalidLons.forEach((lon) => {
        expect(() => GeoValidationMessage_Longitude.parse(lon)).toThrow();
      });
    });
  });

  describe("Accuracy", () => {
    test("valid number does not error", () => {
      validAccWithUndefined.forEach((acc) => {
        expect(() => GeoValidationMessage_Accuracy.parse(acc)).not.toThrow();
      });
    });

    test("valid string-number does not error on safeParse", () => {
      validAcc.forEach((acc) => {
        expect(() =>
          z
            .preprocess((x) => Number(x), GeoValidationMessage_Accuracy)
            .parse(JSON.stringify(acc)),
        ).not.toThrow();
      });
    });

    test("invalid errors", () => {
      invalidAcc.forEach((acc) => {
        expect(() => GeoValidationMessage_Accuracy.parse(acc)).toThrow();
      });
    });
  });

  describe("Message", () => {
    test("valid message does not error", () => {
      validNonStrictMessages.forEach((msg) => {
        expect(() => GeoValidationMessageValidator.parse(msg)).not.toThrow();
      });
    });

    test("invalid errors when strict", () => {
      invalidStrictMessages.forEach((msg) => {
        expect(() => GeoValidationMessageValidator.strict().parse(msg)).toThrow();
      });
    });

    test("invalid errors", () => {
      invalidMessages.forEach((msg) => {
        expect(() => GeoValidationMessageValidator.parse(msg)).toThrow();
      });
    });
  });
});
