import React, { FC, useEffect, useState } from 'react';
import { useSelector } from 'react-redux';

import { AgTable } from '@navi/web-ui/lib/components';
import { Avatar, Typography } from '@navi/web-ui/lib/primitives';

import { selectAdminBoardData } from '@store/Admin/AdminSlice';
import styles from '@src/Pages/KleosDashboard/Partials/LeaderBoardAndGive.module.scss';
import { adminBoardData } from '@src/Pages/KleosDashboard/types';

import './AdminBoardTable.scss';

type rowDataProp = {
  employees: string;
  kudosGiven: number;
  kudosReceived: number;
  profileUrl: string;
};

const AdminBoardTable: FC = () => {
  const myData = useSelector(selectAdminBoardData);
  const [rowCount, setRowCount] = useState<number>(0);
  useEffect(() => {
    setRowCount(myData?.length || 0);
  }, [myData]);

  const returnAvatar = (rowData: rowDataProp): JSX.Element => {
    if (!rowData.employees.length) return <>-</>;
    return (
      <div className={styles['leader-board-team-meta']}>
        {rowData.profileUrl.length > 0 ? (
          <Avatar size={24} isImage src={rowData.profileUrl} />
        ) : (
          <Avatar size={19} sizeVariant="sm">
            {rowData.employees?.charAt(0)?.toLocaleUpperCase()}
          </Avatar>
        )}
        <Typography variant="p3" color={'#585757'}>
          {rowData.employees}
        </Typography>
      </div>
    );
  };

  const columnsDefs = [
    {
      headerName: `Employees (${rowCount})`,
      suppressMovable: true,
      cellRenderer: ({ data }): JSX.Element => {
        return returnAvatar(data);
      },
      flex: 1,
    },
    {
      headerName: 'Kudos Given',
      field: 'kudosGiven',
      suppressMovable: true,
      cellStyle: { marginLeft: '20px' },
      flex: 1,
    },
    {
      headerName: 'Kudos Received',
      field: 'kudosReceived',
      suppressMovable: true,
      cellStyle: { marginLeft: '20px' },
      flex: 1,
    },
  ];

  const adminBoardInfo =
    myData?.map((item: adminBoardData) => {
      return {
        employees: item?.user?.userName,
        kudosGiven: item.kleosGiven,
        kudosReceived: item?.kleosReceived,
        profileUrl: item?.user?.profileUrl,
      };
    }) || [];

  const onGridSizeChanged = (params): void => {
    if (params.api.sizeColumnsToFit) params.api.sizeColumnsToFit();
  };

  return (
    <div className={'table-element'}>
      <AgTable
        theme="alpine"
        sizeColumnsToFit
        domLayout={'normal'}
        headerHeight={54}
        rowHeight={54}
        columnDefs={columnsDefs}
        rowData={adminBoardInfo}
        onGridSizeChanged={onGridSizeChanged}
        suppressColumnMoveAnimation
      />
    </div>
  );
};

export default AdminBoardTable;
