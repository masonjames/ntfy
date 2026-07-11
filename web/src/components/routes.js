import config from "../app/config";
import { shortUrl } from "../app/utils";

const routes = {
  app: config.app_root,
  login: "/login",
  signup: "/signup",
  account: "/account",
  settings: "/settings",
  passwordResetRequest: "/reset-password",
  passwordReset: "/account/password/reset/:token",
  emailVerify: "/account/email/verify/:token",
  subscription: "/:topic",
  subscriptionExternal: "/:baseUrl/:topic",
  forSubscription: (subscription) => {
    if (subscription.baseUrl !== config.base_url) {
      return `/${shortUrl(subscription.baseUrl)}/${subscription.topic}`;
    }
    return `/${subscription.topic}`;
  },
};

const dynamicPathPrefix = (route) => route.slice(0, route.indexOf("/:"));

export const isPublicAuthPath = (pathname) =>
  pathname === routes.login ||
  pathname === routes.passwordResetRequest ||
  pathname.startsWith(`${dynamicPathPrefix(routes.passwordReset)}/`) ||
  pathname.startsWith(`${dynamicPathPrefix(routes.emailVerify)}/`);

export default routes;
