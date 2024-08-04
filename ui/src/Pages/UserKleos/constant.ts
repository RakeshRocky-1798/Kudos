import { PaginatedDataType } from './types';

export const DEFAULT_PAGE_SIZE = 10;
export const DEFAULT_PAGE_NUMBER = 0;

export enum STATUS {
  IDLE = 'idle',
  LOADING = 'loading',
  ERROR = 'error',
}

export const GET_USER_KUDOS_ENDPOINT = (
  email: string,
  type: PaginatedDataType,
  pageNumber: number,
  pageSize: number,
): string =>
  `${window.config.API_BASE_URL}/getPaginatedInfo/${email}?data_type=${type}&page_number=${pageNumber}&page_size=${pageSize}`;
