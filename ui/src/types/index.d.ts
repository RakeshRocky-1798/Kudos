export {};

declare module '*.module.scss';
interface AppConfig {
  ENV: string;
  AUTH_BASE_URL: string;
  AUTH_CLIENT_ID: string;
  API_BASE_URL: string;
  LEADERBOARD_TYPE: string;
  SLACK_CHANNEL_ID: string;
}

declare global {
  interface Window {
    config: AppConfig;
  }
}
