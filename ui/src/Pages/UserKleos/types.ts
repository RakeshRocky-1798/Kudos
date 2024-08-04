export interface PagesInfo {
  pageNumber: number;
  pageSize: number;
  totalPages: number;
  totalElements: number;
  hasData: boolean;
}

export interface UserData {
  email: string;
  profileUrl: string;
  userName: string;
}
export interface AchievementOption {
  aCreatedAt: string;
  aEmoji: string;
  aType: string;
  user: UserData;
}
export interface receivedTabData {
  id: number;
  message?: string;
  achievementData: AchievementOption;
}

export interface userKudosApiResponse {
  data: Array<receivedTabData>;
  pages: PagesInfo;
  status?: string;
}

// export type SessionsType = receivedTabType;
export type SessionsData = receivedTabData;
export type PaginatedDataType = 'given' | 'received';
