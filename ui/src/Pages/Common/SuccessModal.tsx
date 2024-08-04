import { FC } from 'react';
import { useDispatch, useSelector } from 'react-redux';

import Button from '@navi/web-ui/lib/primitives/Button';
import ModalDialog from '@navi/web-ui/lib/primitives/ModalDialog';
import Typography from '@navi/web-ui/lib/primitives/Typography';
import {
  setGiveKleosModalState,
  setSuccessModalState,
  setIsGiven,
  selectIsGiven,
} from '@store/Dashboard/DashboardSlice';

import SuccessAnimation from './SuccessAnimation';
import styles from './SuccessModal.module.scss';

interface SuccessModalProps {
  isOpen: boolean;
  onClose: () => void;
}

const SuccessModal: FC<SuccessModalProps> = ({ isOpen, onClose }) => {
  const dispatch = useDispatch();
  const isGiven = useSelector(selectIsGiven);

  const handleGiveMoreKleos = (): void => {
    dispatch(setSuccessModalState(false));
    dispatch(setGiveKleosModalState(true));
  };

  const handleDashboardNavigation = (): void => {
    dispatch(setSuccessModalState(false));
    dispatch(setIsGiven(true));
  };

  return (
    <ModalDialog
      customFooter={
        <div className={styles['giveMoreKudos']}>
          {isGiven ? (
            <Button
              type="submit"
              variant="secondary"
              className={styles['give-kleos-button']}
              onClick={onClose}
            >
              Cancel
            </Button>
          ) : (
            <Button
              type="submit"
              variant="primary"
              className={styles['give-kleos-button']}
              onClick={handleDashboardNavigation}
            >
              Continue to the Dashboard
            </Button>
          )}

          <span className={styles['median-span']}></span>
          {isGiven && (
            <Button
              type="submit"
              variant="primary"
              className={styles['give-kleos-button']}
              onClick={handleGiveMoreKleos}
            >
              Give More Kudos
            </Button>
          )}
          {/* <Button
            type="submit"
            variant="primary"
            className={styles['give-kleos-button']}
            onClick={handleGiveMoreKleos}
          >
            Give More Kudos
          </Button> */}
        </div>
      }
      onClose={onClose}
      open={isOpen}
    >
      <SuccessAnimation />
      <div className={styles['headDiv']}>
        <Typography
          variant="h1"
          className={styles['hurray']}
          color={'#1C1C1C'}
          style={{ fontWeight: '600' }}
        >
          Hurrayy!!
        </Typography>
        <Typography
          variant="p3"
          className={styles['kudDesc']}
          color={'#657384'}
          style={{ fontSize: '16px', fontWeight: '400' }}
        >
          Kudos given, keep motivating the team
        </Typography>
      </div>
    </ModalDialog>
  );
};

export default SuccessModal;
