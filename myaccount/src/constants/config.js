import _ from 'lodash';

// environment names
export const ENV_NAME_DEVELOPMENT = 'DEV';
export const ENV_NAME_STAGING = 'STAGE';
export const ENV_NAME_PRODUCTION = 'PROD';

/*eslint-disable */
export const SENTRY_PRODUCTION_KEY = 'https://f4d735c8d9fc44e2a184f6dd858302c4@sentry.io/1496206';
export const SENTRY_STAGING_KEY = 'https://f4d735c8d9fc44e2a184f6dd858302c4@sentry.io/1496206';
/*eslint-enable */

// apex for the various staffjoy environments
export const HTTP_PREFIX = 'http://';
export const HTTPS_PREFIX = 'https://';
export const DEVELOPMENT_APEX = '.staffjoy-v2.local';
export const STAGING_APEX = '.staffjoystaging.com';
export const PRODUCTION_APEX = '.staffjoy.com';

const DEFAULT_REFETCH_INTERVAL = 10;

const REFETCH_INTERVALS = {
  USER: 30,
  WHOAMI: 30,
  DEFAULT: DEFAULT_REFETCH_INTERVAL,
};

export function getRefetchInterval(endpoint) {
  return _.get(REFETCH_INTERVALS, endpoint, DEFAULT_REFETCH_INTERVAL);
}

export const MOMENT_MONTH_YEAR_FORMAT = 'MMMM YYYY';
