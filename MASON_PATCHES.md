# Mason fork patch ledger

These patches are intentionally carried on `masonjames/ntfy` until equivalent
upstream fixes land. Weekly upstream synchronization must preserve their tests;
a conflict opens an issue instead of silently dropping either behavior.

## Password reset confirmation without SMTP

- **Source:** [PR #1 review](https://github.com/masonjames/ntfy/pull/1#discussion_r3564996144)
- **Behavior:** `POST /v1/account/password/reset` requires the user manager, but
  not a configured mailer. This keeps `ntfy user reset-pass` links usable when
  the CLI prints a token for an installation without SMTP.
- **Regression test:** `TestAccount_PasswordReset_ConfirmWithoutSMTP`.
- **Removal criterion:** upstream accepts an equivalent route guard and test.

## Logged-out magic-link routes with required login

- **Source:** [PR #1 review](https://github.com/masonjames/ntfy/pull/1#discussion_r3564996145)
- **Behavior:** the login redirect permits password-reset request/confirmation
  and email-verification routes for logged-out visitors. Protected app routes
  still redirect to `/login`.
- **Regression test:** `web/src/components/routes.test.js`.
- **Removal criterion:** upstream accepts equivalent route matching and tests.
