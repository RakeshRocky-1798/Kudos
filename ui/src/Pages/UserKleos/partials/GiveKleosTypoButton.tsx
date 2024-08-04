import { FC, useState } from 'react';
import { Typography, Button } from '@navi/web-ui/lib/primitives';
import { useDispatch } from 'react-redux';

import GiveKleosModal from '@src/Pages/UserKleos/partials/GiveKleosModal';

import style from './GiveKleosTypoButton.module.scss';
import {
  selectGiveKleosModalState,
  setGiveKleosModalState,
} from '@store/Dashboard/DashboardSlice';
import { useSelector } from 'react-redux';

const GiveKleosTypoButton: FC = () => {
  const dispatch = useDispatch();
  const modalShow = useSelector(selectGiveKleosModalState);
  const toggleModal = (): void => {
    dispatch(setGiveKleosModalState(!modalShow));
  };

  return (
    <div className={style['main-wrapper']}>
      <div className={style['typo-wrapper']}>
        <Typography
          variant="h3"
          color={'#1C1C1C'}
          style={{ fontSize: '16px' }}
          className={style['typo-inner-one']}
        >
          Appreciate your fellow Navi_ites
        </Typography>
        <Typography
          variant="p4"
          color={'#585757'}
          style={{ fontSize: '12px' }}
          className={style['typo-inner-two']}
        >
          Take a moment to recognise and appreciate your teammates
        </Typography>
      </div>
      <div className={style['button-wrapper']}>
        <Button variant="primary" size="medium" onClick={toggleModal}>
          Give Kudos
        </Button>
        <GiveKleosModal onClose={toggleModal} open={modalShow} />
      </div>
    </div>
  );
};

export default GiveKleosTypoButton;
