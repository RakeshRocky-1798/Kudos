import { FC } from 'react';
import { useSelector } from 'react-redux';
import cx from 'classnames';

import { AgTable } from '@navi/web-ui/lib/components';
import { Avatar, Typography } from '@navi/web-ui/lib/primitives';

import LockIcon from '@src/assets/icons/lock.svg';
import { selectLeaderBoardData } from '@src/store/Dashboard/DashboardSlice';
import { getFromLocalStorage } from '@src/utils/storage';

import { returnLeaderBoardPositionEmoji } from '../constants';
import { leaderBoardData } from '../types';
import styles from './LeaderBoardAndGive.module.scss';
import './LeaderBoardTable.scss';

type rowDataProp = {
  position: number;
  count: number;
  teamMember: string;
  teamMemberProfileUrl: string;
  userEmail: string;
};

const LeaderBoardTable: FC = () => {
  const currentDate = new Date();
  const currentMonth = currentDate.toLocaleString('default', { month: 'long' });
  const currentYear = currentDate.getFullYear();

  const leaderBoardData: leaderBoardData[] =
    useSelector(selectLeaderBoardData) || [];
  const userName: string = getFromLocalStorage('user-name') || '';
  const loggedInUserProfileUrl = getFromLocalStorage('user-profile-url') || '';
  const email = getFromLocalStorage('email-id') || '';
  const checkIsLoggedInUser = (teamMember: string): boolean =>
    teamMember === email;

  const returnAvatar = (rowData: rowDataProp): JSX.Element => {
    const isCurrentUser = checkIsLoggedInUser(rowData.userEmail);

    if (!rowData?.teamMember.length) {
      return <>-</>;
    }

    return (
      <div className={styles['leader-board-team-meta']}>
        {rowData.teamMemberProfileUrl.length > 0 ? (
          <Avatar size={24} isImage src={rowData.teamMemberProfileUrl} />
        ) : (
          <Avatar size={24} sizeVariant="sm">
            {rowData.teamMember?.charAt(0)?.toLocaleUpperCase()}
          </Avatar>
        )}
        <Typography
          variant="p3"
          className={
            isCurrentUser
              ? styles['logged-in-user']
              : styles['leader-board-user']
          }
        >
          {isCurrentUser ? 'You' : rowData.teamMember}
        </Typography>
      </div>
    );
  };

  const leaderBoardInfo = leaderBoardData.map((item: leaderBoardData) => {
    return {
      position: item.rank,
      count: item.totalCount,
      teamMember: item?.userMeta?.userName,
      teamMemberProfileUrl: item?.userMeta?.profileUrl,
      userEmail: item?.userMeta?.email,
    };
  });

  const columnsDef = [
    {
      headerName: 'Position',
      field: 'position',
      suppressMovable: true,
      cellRenderer: ({ data }): JSX.Element => {
        return (
          <div
            className={cx(styles['leader-board-team-meta'], {
              [styles['logged-in-user']]: checkIsLoggedInUser(data.userEmail),
            })}
          >
            <Typography variant="p3"># {data.position}</Typography>
            <span>{returnLeaderBoardPositionEmoji(data.position)}</span>
          </div>
        );
      },
      cellStyle: { marginLeft: '8px' },
      flex: 0.5,
    },
    {
      // TODO: Check in store for the Kudos given / received config then update the header name
      headerName: 'Kudos Given',
      field: 'count',
      suppressMovable: true,
      cellStyle: {
        marginLeft: '20px',
      },
      flex: 0.5,
      cellRenderer: ({ data }): JSX.Element => (
        <Typography
          variant="p3"
          className={cx({
            [styles['logged-in-user']]: checkIsLoggedInUser(data.userEmail),
          })}
        >
          {data.count || 0}
        </Typography>
      ),
    },
    {
      headerName: 'Name',
      suppressMovable: true,
      cellRenderer: ({ data }): JSX.Element => {
        return returnAvatar(data);
      },
      flex: 1,
    },
  ];

  const onGridSizeChanged = (params): void => {
    if (params.api.sizeColumnsToFit) params.api.sizeColumnsToFit();
  };

  const returnRowStyle = (params): { background: string } => {
    if (checkIsLoggedInUser(params?.data?.userEmail)) {
      return {
        background: 'var(--navi-color-blue-bg)',
      };
    }
    return {
      background: 'transparent',
    };
  };

  if (!leaderBoardInfo.length) {
    return (
      <section>
        <Typography
          variant="p3"
          className={styles['leader-board-table-header']}
        >
          Appreciation Champions Board
        </Typography>
        <div className={styles['empty-leader-board']}>
          <img src={LockIcon} alt="lock" />
          <Typography variant="p3" className={styles['empty-leader-board--ms']}>
            Patience, mateys!
            <Typography variant="p3">
              This treasure chest pops open with a flood of Kudos!
            </Typography>
          </Typography>
        </div>
      </section>
    );
  }

  return (
    <div className="leader-board-table">
      <Typography variant="p3" className={styles['leader-board-table-header']}>
        Appreciation Champions Board
      </Typography>
      <Typography variant="p3" className={styles['leader-board-table-header']}>
        {currentMonth} {currentYear}
      </Typography>
      <AgTable
        theme="alpine"
        sizeColumnsToFit
        domLayout="normal"
        headerHeight={54}
        rowHeight={58}
        columnDefs={columnsDef}
        rowData={leaderBoardInfo}
        getRowStyle={returnRowStyle}
        onGridSizeChanged={onGridSizeChanged}
        suppressColumnMoveAnimation
      />
    </div>
  );
};

export default LeaderBoardTable;
