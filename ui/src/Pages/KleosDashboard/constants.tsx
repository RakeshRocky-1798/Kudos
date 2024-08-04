import { TagColors } from '@navi/web-ui/lib/primitives/Tag/types';

export const FETCH_MY_DATA_URL = (): string => {
  return `${window?.config?.API_BASE_URL}/kleosclient/getGroup/U05DY7YQMR9`;
};

export const FETCH_USER_KLEOS_DATA = (emailId: string): string => {
  return `${window?.config?.API_BASE_URL}/getKleos/${emailId}`;
};

export const FETCH_DASHBOARD_LEADERBOARD_DATA = (
  isReceived: boolean,
  emailId: string,
): string => {
  return `${window?.config?.API_BASE_URL}/leaderboard?isReceived=${isReceived}&userid=${emailId}`;
};

export const FETCH_FIRST_PAGE_DATA = (emailId: string): string => {
  return `${window?.config?.API_BASE_URL}/getKleosReceived/${emailId}`;
};

export const FETCH_ALL_USERS = (emailId: string): string => {
  return `${window?.config?.API_BASE_URL}/getAllUsers?userid=${emailId}`;
};

export const FETCH_ALL_ACHIEVEMENTS = (): string => {
  return `${window?.config?.API_BASE_URL}/achievement`;
};

export const GIVE_KLEOS = (): string => {
  return `${window?.config?.API_BASE_URL}/giveKleos`;
};

export const FETCH_MY_KLEOS_COUNT_URL = (emailId: string): string => {
  return `${window?.config?.API_BASE_URL}/getKleos/${emailId}`;
};

export const ADMIN_DATA = (): string => {
  return `${window?.config?.API_BASE_URL}/getAdminData`;
};

export const DOWNLOAD_XLS = (emailId: string): string => {
  return `${window?.config?.API_BASE_URL}/getAdminData/xls/${emailId}`;
};

export const FETCH_TAB_CONTENT = (
  emailId: string,
  dataType: string,
  pageSize: number,
  pageNumber: number,
): string => {
  return `${window?.config?.API_BASE_URL}/getPaginatedInfo/${emailId}?data_type=${dataType}&page_size=${pageSize}&page_number=${pageNumber}`;
};

export const returnAchievementIcon = (aEmoji: string): string | null => {
  switch (aEmoji) {
    case 'star':
      return '⭐';
    case 'muscle':
      return '💪';
    case 'fire':
      return '🔥';
    case 'leaves':
      return '🍃';
    case 'trophy':
      return '🏆';
    case 'thumbsup':
      return '👍';
    case 'sports_medal':
      return '🏅';
    case 'rocket':
      return '🚀';
    default:
      return null;
  }
};

export const returnAchieveTagVariantColor = (aEmoji: string): TagColors => {
  switch (aEmoji) {
    case 'star':
    case 'trophy':
    case 'thumbsup':
      return 'yellow';
    case 'muscle':
    case 'rocket':
      return 'green';
    case 'leaves':
    case 'sports_medal':
      return 'blue';
    case 'fire':
      return 'red';
    default:
      return 'primary';
  }
};

export const returnLeaderBoardPositionEmoji = (
  position: number,
): React.ReactNode | null => {
  switch (position) {
    case 1:
      return '🥇';
    case 2:
      return '🥈';
    case 3:
      return '🥉';
    default:
      return null;
  }
};

export const USER_PERMISSION = {
  REQUEST_ADMIN_BOARD_READ_PERMISSION: 'user.admin_panel.read',
  REQUEST_ADMIN_BOARD_WRITE_PERMISSION: 'user.admin_panel.write',
};
export const ROLES = {
  ADMIN: 'admin',
  ALLUSERS: 'allusers',
};
