import { DropdownOptionType } from '@src/Pages/KleosDashboard/types';

export interface UserData {
  data: Array<DropdownOptionType>;
  error: Object | unknown;
  status: string;
}

export interface AchievementData {
  data: Array<DropdownOptionType>;
  error: Object | unknown;
  status: string;
}
