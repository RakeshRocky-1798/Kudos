export type DropdownOptionType = {
  label: string;
  value: string;
  aEmoji: string;
};

export interface AchievementData {
  achievementName: string;
  count: number;
  emoji: string;
}

export interface recentRecognitionData {
  message: string;
  achievement: {
    aLabel: string;
    aEmoji: string;
    aFrom: {
      email: string;
      userName: string;
      profileUrl: string;
    };
  };
}

export interface leaderBoardData {
  rank: number;
  totalCount: number;
  userMeta: {
    email: string;
    userName: string;
    profileUrl: string;
  };
}

export interface adminBoardData {
  kleosGiven: number;
  kleosReceived: number;
  user: {
    email: string;
    userName: string;
    profileUrl: string;
  };
}

export interface DashboardData {
  myData: {
    email: string;
    userNumber: string;
    profileUrl: string;
  };
  kleosMetrics: {
    givenCount: number;
    receivedCount: number;
  };
  achievementDropdown: Array<DropdownOptionType>;
  totalAchievement: Array<AchievementData>;
  recentRecognition: Array<recentRecognitionData>;
}

export interface UserKleosData {
  data: {
    myData: {
      email: string;
      userNumber: string;
      profileUrl: string;
    };
    kleosMetrics: {
      givenCount: number;
      receivedCount: number;
    };
    achievementDropdown: Array<DropdownOptionType>;
    totalAchievement: Array<AchievementData>;
    recentRecognition: Array<recentRecognitionData>;
  };
  error: Object | unknown;
  status: string;
}

export interface DashboardLeaderBoardData {
  data: {
    leaderBoardData: Array<leaderBoardData>;
  };
  error: Object | unknown;
  status: string;
}

export interface AdminPortalBoardData {
  data: {
    adminAllUser: Array<adminBoardData>;
  };
  error: Object | unknown;
  status: string;
}
