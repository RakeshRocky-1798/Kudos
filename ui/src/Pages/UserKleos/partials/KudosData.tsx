import { FC } from 'react';
import { useSelector } from 'react-redux';

import { Typography, Avatar, Tag } from '@navi/web-ui/lib/primitives';
import {
  selectReceivedTabData,
  selectGivenTabData,
  selectCurrentTab,
} from '@src/store/UserKudos/UserSlice';
import {
  returnAchieveTagVariantColor,
  returnAchievementIcon,
} from '@src/Pages/KleosDashboard/constants';
import style from './UserKudosTabs.module.scss';
import { receivedTabData } from '../types';

const UserBlock: FC<{ item: receivedTabData }> = ({ item }) => {
  return (
    <>
      <div className={style['kudos-user']}>
        {item?.achievementData?.user?.profileUrl.length > 0 ? (
          <Avatar
            src={item?.achievementData?.user?.profileUrl}
            alt="User Avatar"
            sizeVariant="sm"
            isImage
          />
        ) : (
          <Avatar sizeVariant="sm">
            {item?.achievementData?.user?.userName.charAt(0)?.toUpperCase()}
          </Avatar>
        )}
        <Typography variant="p3">
          {item?.achievementData?.user?.userName}
        </Typography>
        <Tag
          label={item?.achievementData?.aType}
          variant="transparent"
          color={returnAchieveTagVariantColor(item?.achievementData?.aEmoji)}
          endAdornment={returnAchievementIcon(item?.achievementData?.aEmoji)}
        />
      </div>
    </>
  );
};

const TimeDetails: FC<{ item: receivedTabData }> = ({ item }) => {
  return (
    <Typography variant="p5" color="var(--navi-color-gray-c3)">
      {item?.achievementData?.aCreatedAt && (
        <>{item?.achievementData?.aCreatedAt?.replace(',', ' |')}</>
      )}
    </Typography>
  );
};

const KudosData: FC = () => {
  const receivedTabData = useSelector(selectReceivedTabData);
  const givenTabData = useSelector(selectGivenTabData);
  const currentTab = useSelector(selectCurrentTab);
  const data =
    currentTab === 'received' ? receivedTabData.data : givenTabData.data;

  const returnCurrentTabContent = (item: receivedTabData): JSX.Element => {
    return (
      <div className={style['kudos-given-tab-item']}>
        <div key={item?.id}>
          <div className={style['kudos-tab-header-wrapper']}>
            <UserBlock item={item} />
            <TimeDetails item={item} />
          </div>
          <section className={style['kudos-user-achievement-wrapper']}>
            <Typography variant="p3">{item?.message}</Typography>
          </section>
        </div>
        <div className={style['divider']} />
      </div>
    );
  };

  return (
    <div className={style['tab-item-wrapper']}>
      {data?.map((item: receivedTabData) => returnCurrentTabContent(item))}
    </div>
  );
};

export default KudosData;
