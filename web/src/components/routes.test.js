import { describe, expect, it } from "vitest";
import { isPublicAuthPath } from "./routes";

describe("isPublicAuthPath", () => {
  it.each(["/login", "/reset-password", "/account/password/reset/reset-token", "/account/email/verify/verify-token"])(
    "allows logged-out authentication route %s",
    (pathname) => {
      expect(isPublicAuthPath(pathname)).toBe(true);
    },
  );

  it.each(["/", "/account", "/settings", "/account/password/reset", "/account/email/verify"])(
    "keeps protected or incomplete route %s behind login",
    (pathname) => {
      expect(isPublicAuthPath(pathname)).toBe(false);
    },
  );
});
