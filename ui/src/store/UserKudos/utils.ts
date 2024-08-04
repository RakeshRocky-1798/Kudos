import {
  userKudosApiResponse,
  receivedTabData,
} from '@src/Pages/UserKleos/types';

const isDateFormatProper = (date: string): boolean => {
  const tempDate: Date = new Date(date);
  return tempDate?.toString() !== 'Invalid Date';
};

export const transformCreatedAt = (createdAt: string): string => {
  if (createdAt?.length) {
    const date = new Date(createdAt);
    // Adjusting for IST (GMT+5:30)
    date.setHours(date.getHours() + 5);
    date.setMinutes(date.getMinutes() + 30);

    const year = date.getFullYear();
    const month = date.getMonth() + 1; // Month is zero-based, so we add 1
    const monthString = month < 10 ? `0${month}` : `${month}`;
    const day = date.getDate();
    const dayString = day < 10 ? `0${day}` : `${day}`;
    let hours = date.getHours();
    const meridiem = hours >= 12 ? 'PM' : 'AM';

    // Convert to 12-hour format
    hours = hours % 12 || 12; // Convert hours to 12-hour format
    const finalHours = hours < 10 ? `0${hours}` : `${hours}`;

    const minutes = date.getMinutes();
    const finalMinutes = minutes < 10 ? `0${minutes}` : `${minutes}`;

    return `${dayString}-${monthString}-${year}, ${finalHours}:${finalMinutes} ${meridiem}`;
  }
  return '';
};

export const transformUserKudosData = (
  existingData: receivedTabData[],
  ApiResponseData: userKudosApiResponse,
): userKudosApiResponse => {
  const tempArr: receivedTabData[] = [];
  const temp = {};
  const finalData: userKudosApiResponse = {
    data: [],
    pages: {
      pageNumber: 0,
      pageSize: 10,
      totalPages: 0,
      totalElements: 0,
      hasData: false,
    },
  };

  if (existingData?.length) {
    existingData.forEach((data: receivedTabData) => {
      tempArr.push(data);
    });
  }

  if (ApiResponseData?.data?.length) {
    ApiResponseData.data.forEach((data: receivedTabData) => {
      tempArr.push(data);
    });
  }

  tempArr?.forEach((data: receivedTabData) => {
    if (!temp[data?.id]) {
      const finalPayload = {
        id: data?.id,
        message: data?.message,
        achievementData: {
          aCreatedAt: isDateFormatProper(data.achievementData.aCreatedAt)
            ? transformCreatedAt(data.achievementData.aCreatedAt)
            : data.achievementData.aCreatedAt,
          aEmoji: data?.achievementData?.aEmoji,
          aType: data?.achievementData?.aType,
          user: {
            email: data?.achievementData?.user?.email,
            profileUrl: data?.achievementData?.user?.profileUrl,
            userName: data?.achievementData?.user?.userName,
          },
        },
      };
      finalData.data.push(finalPayload);
      temp[data.id] = true;
    }
  });

  finalData.pages = ApiResponseData.pages;
  return finalData;
};
