import { FC } from 'react';
import { useDispatch, useSelector } from 'react-redux';

import { Tabs, TabItem } from '@navi/web-ui/lib/components';
import { TabItemKey } from '@navi/web-ui/lib/components/Tabs/types';

import RaiseHand from '@src/assets/icons/raiseHand.svg';
import Clap from '@src/assets/icons/clap.svg';
import {
  selectCurrentTab,
  selectGivenTabData,
  selectReceivedTabData,
  setCurrentTab,
} from '@src/store/UserKudos/UserSlice';
import { selectDashboardData } from '@src/store/Dashboard/DashboardSlice';
import useUserKudosApi from '@src/Pages/UserKleos/useUserKudosApi';

import KudosGivenTab from './KudosGivenTab';
import KudosReceivedTab from './KudosReceivedTab';
import { PaginatedDataType } from '../types';
import style from './UserKudosTabs.module.scss';

const UserKudosTabs: FC = () => {
  const currentTab = useSelector(selectCurrentTab);
  const userData = useSelector(selectDashboardData);
  const givenTabData = useSelector(selectGivenTabData);
  const receivedTabData = useSelector(selectReceivedTabData);
  const dispatch = useDispatch();
  const { fetchUserKudosData } = useUserKudosApi();

  const returnTabItemCounter = (key: PaginatedDataType): number => {
    if (key === 'given') {
      return (
        givenTabData.pages.totalElements || userData.kleosMetrics.givenCount
      );
    }
    return (
      receivedTabData.pages.totalElements || userData.kleosMetrics.receivedCount
    );
  };

  const returnTabItemAdornment = (key: PaginatedDataType): JSX.Element => {
    if (key === 'given') {
      return <img src={Clap} alt="Clap Hand Emoji" />;
    }
    return <img src={RaiseHand} alt="Raise Hand Emoji" />;
  };

  const handleTabChange = (key: TabItemKey): void => {
    dispatch(setCurrentTab(key as PaginatedDataType));
    fetchUserKudosData(0, key as PaginatedDataType, false);
  };

  return (
    <Tabs
      selectedTabKey={currentTab}
      onTabChange={handleTabChange}
      tabsClassName={style['tab-wrapper']}
      contentClassName={style['tab-content-wrapper']}
    >
      <TabItem
        key={'given'}
        label={'Given'}
        counter={returnTabItemCounter('given')}
        startAdornment={returnTabItemAdornment('given')}
      >
        <KudosGivenTab />
      </TabItem>
      <TabItem
        key={'received'}
        label={'Received'}
        counter={returnTabItemCounter('received')}
        startAdornment={returnTabItemAdornment('received')}
      >
        <KudosReceivedTab />
      </TabItem>
    </Tabs>
  );
};

export default UserKudosTabs;
