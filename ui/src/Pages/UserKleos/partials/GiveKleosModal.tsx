import { FC, useState } from 'react';
import { ModalDialog, Typography } from '@navi/web-ui/lib/primitives';
import GiveKleosForm from '@src/Pages/Common/GiveKleosForm';
import { useSelector } from 'react-redux';
import {
  selectAllAchievementData,
  selectAllUserData,
  selectCurrentUser,
} from '@store/Dashboard/DashboardSlice';

import styles from './GiveKleosModal.module.scss';
import { returnAchievementIcon } from '@src/Pages/KleosDashboard/constants';

interface GiveKleosModalProps {
  open: boolean;
  onClose: () => void;
}

const GiveKleosModal: FC<GiveKleosModalProps> = ({ open, onClose }) => {
  const userId = useSelector(selectCurrentUser);
  const allUsers = useSelector(selectAllUserData).data;
  const allAchievements = useSelector(selectAllAchievementData).data;

  if (allUsers.length !== 0 || allAchievements.length !== 0) {
    return (
      <ModalDialog header={'Give Kudos'} onClose={onClose} open={open}>
        <Typography variant="p3" className={styles['giveKudosDesc']}>
          Give them a high-five so high it reaches the moon!
        </Typography>
        <GiveKleosForm
          toOptions={allUsers}
          achievementOptions={allAchievements}
          userId={userId}
        />
      </ModalDialog>
    );
  }
  return null;
};

export default GiveKleosModal;
