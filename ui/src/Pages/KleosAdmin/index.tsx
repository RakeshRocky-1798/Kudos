import { FC, useEffect } from 'react';

import Button from '@navi/web-ui/lib/primitives/Button';
import { toast } from '@navi/web-ui/lib/primitives/Toast';

import { ApiService } from '@src/service/api';
import ErrorBoundary from '@src/components/ErrorBoundary/ErrorBoundary';
import FallbackComponent from '@src/components/Fallback';
import AdminBoardTable from '@src/Pages/KleosAdmin/AdminBoardTable';
import styles from '@src/Pages/KleosAdmin/AdminBoard.module.scss';
import useKleosAdminApis from '@src/Pages/KleosAdmin/adminDataAPI';
import { DOWNLOAD_XLS } from '@src/Pages/KleosDashboard/constants';
import { getFromLocalStorage } from '@src/service/storage';

const AdminBoard: FC = () => {
  const { fetchAndSetAdminBoardData } = useKleosAdminApis();

  const fetchAllData = (): void => {
    fetchAndSetAdminBoardData();
  };

  useEffect(() => {
    fetchAllData();
  }, []);

  const downloadCsv = (): void => {
    const emailId: string = getFromLocalStorage('email-id') || '';
    ApiService.get(DOWNLOAD_XLS(emailId))
      .then(res => {
        if (res.data?.status === 200) {
          toast.success('Report will be shared over slack 😇');
        }
      })
      .catch(err => {
        toast.error('Error in generating report');
      });
  };

  const returnContent = (): JSX.Element => {
    return (
      <>
        <div className={styles['main-div']}>
          <div className={styles['child-div']}>
            <AdminBoardTable />
            <div className={styles['button-wrapper']}>
              <Button
                type="submit"
                variant="primary"
                className={styles['give-kleos-button']}
                onClick={downloadCsv}
                disabled={false}
              >
                Generate Report
              </Button>
            </div>
          </div>
        </div>
      </>
    );
  };

  return (
    <ErrorBoundary fallbackComponent={<FallbackComponent />}>
      {returnContent()}
    </ErrorBoundary>
  );
};

export default AdminBoard;
