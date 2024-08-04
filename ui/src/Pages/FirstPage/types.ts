import { DropdownOptionType } from '@src/Pages/KleosDashboard/types';

export interface HomeData {
  data: {
    userEmail: string;
    givenCount: string;
    receivedCount: string;
  };
  error: Object | unknown;
  status: string;
}
