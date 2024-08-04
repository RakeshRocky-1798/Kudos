import React, { FC, ReactNode, useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useDispatch, useSelector } from 'react-redux';
import { Select } from 'antd';
import {
  Button,
  Switch,
  TextArea,
  Typography,
} from '@navi/web-ui/lib/primitives';
import { filterOption } from '@utils/helper';
import { ApiService } from '@src/service/api';
import { GIVE_KLEOS } from '@src/Pages/KleosDashboard/constants';

import {
  selectSuccessModalState,
  setGiveKleosModalState,
  setIsGiven,
  setSuccessModalState,
  selectIsGiven,
} from '@store/Dashboard/DashboardSlice';

import { returnAchievementIcon } from '@src/Pages/KleosDashboard/constants';

import SuccessModal from './SuccessModal';
import styles from './GiveKleosForm.module.scss';
import { toast } from '@navi/web-ui/lib/primitives/Toast';
import useKleosDashboardApi from '@src/Pages/KleosDashboard/useKleosDashboardApi';
import useUserKudosApi from '@src/Pages/UserKleos/useUserKudosApi';
import { InfoIcon } from '@navi/web-ui/lib/icons';
import { getFromLocalStorage } from '@src/service/storage';

type OptionProps = {
  label: string;
  value: string;
  aEmoji: string;
};

export type DropdownOptionType = {
  label: string;
  value: string;
  aEmoji: string;
};

interface GiveKleosFormProps {
  toOptions: OptionProps[];
  achievementOptions: DropdownOptionType[]; // Change here
  userId: string;
}
const GiveKleosForm: FC<GiveKleosFormProps> = ({
  toOptions,
  achievementOptions,
  userId,
}) => {
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const [messages, setMessage] = useState('');
  const [to, setTo] = useState('');
  const [achievement, setAchievement] = useState('');
  const isSuccessModalOpen = useSelector(selectSuccessModalState);
  const { startFetchingDashboardData } = useKleosDashboardApi();
  const { fetchUserKudosData } = useUserKudosApi();
  const [sendToSlack, setSendToSlack] = useState(true);
  const [buttonEnabled, setButtonEnabled] = useState(false);
  const AchievementWithIcons = achievementOptions.map(option => ({
    ...option,
    label: `${returnAchievementIcon(option.aEmoji)} ${option.label}`,
  }));
  const isGiven = useSelector(selectIsGiven);

  useEffect(() => {
    const isFieldsFilled = to !== '' && achievement !== '' && messages !== '';
    setButtonEnabled(isFieldsFilled);
  }, [to, achievement, messages]);

  const handleAchievementSelect = (option: string): void => {
    setAchievement(option);
  };

  const handleToSelect = (option: string): void => {
    setTo(option);
  };

  const handleTextAreaChange = (
    e: React.ChangeEvent<HTMLTextAreaElement>,
  ): void => {
    setMessage(e.target.value);
  };

  const handleSwitch = (): void => {
    setSendToSlack(!sendToSlack);
  };

  const handleSuccessResponse = (): void => {
    dispatch(setGiveKleosModalState(false));
    setTo('');
    setAchievement('');
    setMessage('');
    setSendToSlack(false);
    dispatch(setSuccessModalState(true));
  };

  const handleModalClose = (): void => {
    dispatch(setSuccessModalState(false));
    const location = window.location.pathname;
    dispatch(setIsGiven(true));

    if (location.includes('user-kudos')) {
      fetchUserKudosData(0, 'given', false);
    } else {
      startFetchingDashboardData();
    }
  };

  const sendKleos = (): void => {
    const body = {
      from: userId,
      to: to,
      message: messages,
      achievement: achievement,
      needSlack: sendToSlack,
    };

    if (to === '' || achievement === '' || messages === '') {
      toast.error('Please fill all the fields');
    } else {
      ApiService.post(GIVE_KLEOS(), body)
        .then(res => {
          if (res.data?.status === 200) {
            handleSuccessResponse();
          }
        })
        .catch(err => {
          const finalMessage =
            err.response?.data?.message || 'Something went wrong';
          toast.error(finalMessage);
        });
    }
  };

  const customStyle = {
    zIndex: 1000000000000000,
  };

  if (isSuccessModalOpen) {
    return (
      <div>
        <SuccessModal isOpen={isSuccessModalOpen} onClose={handleModalClose} />
      </div>
    );
  }

  return (
    <div className={styles['form-div']}>
      <div className={styles['form-element']}>
        <Typography
          variant={'p2'}
          color={'#585757'}
          style={{ fontSize: '13px', fontWeight: '400' }}
        >
          Kudos to
          <span style={{ color: 'var(--navi-required-icon-default-color)' }}>
            {' '}
            *
          </span>
        </Typography>
        <Select
          className={styles['select-box']}
          placeholder="Enter Name to search"
          optionFilterProp="children"
          showSearch={true}
          options={toOptions}
          onChange={handleToSelect}
          mode="multiple"
          filterOption={filterOption}
          dropdownStyle={customStyle}
          suffixIcon={null}
        />
      </div>
      <div className={styles['form-element']}>
        <Typography
          variant={'p2'}
          color={'#585757'}
          style={{ fontSize: '13px', fontWeight: '400' }}
        >
          Appreciation
          <span style={{ color: 'var(--navi-required-icon-default-color)' }}>
            {' '}
            *
          </span>
        </Typography>

        <Select
          className={styles['select-box']}
          placeholder="Select a appreciation type"
          optionFilterProp="children"
          options={AchievementWithIcons}
          onChange={handleAchievementSelect}
          filterOption={filterOption}
          dropdownStyle={customStyle}
        />
      </div>
      <div className={styles['form-element']}>
        <Typography
          variant={'p2'}
          color={'#585757'}
          style={{ fontSize: '13px', fontWeight: '400' }}
        >
          Send a nice note
          <span style={{ color: 'var(--navi-required-icon-default-color)' }}>
            {' '}
            *
          </span>
        </Typography>
        <TextArea
          placeholder="Describe their actions and how they impacted you or the business."
          fullWidth={true}
          onChange={handleTextAreaChange}
          style={{
            width: '10%',
          }}
        />
      </div>
      <div className={styles['send-slack']}>
        <div className={styles['send-slack-btn']}>
          <Switch
            onChange={handleSwitch}
            checked={sendToSlack}
            crossOrigin={undefined}
          />
        </div>
        <Typography variant={'p3'} color={'#657384'}>
          Post in{' '}
          <a
            style={{ textDecoration: 'none' }}
            href={`https://go-navi.slack.com/archives/${window.config.SLACK_CHANNEL_ID}`}
            target="_blank"
            rel="noreferrer"
          >
            <Typography
              as="span"
              variant="h4"
              color="var(--navi-color-blue-base)"
            >
              #kudos
            </Typography>
          </a>{' '}
          slack channel
        </Typography>
      </div>
      <div className={styles['info-div']}>
        <InfoIcon className={styles['info-icon']} />
        <Typography variant={'p4'} color={'#657384'}>
          Kudos will also be shared to the recipient&apos;s manager{' '}
        </Typography>
      </div>
      <div className={styles['give-kleos-button']}>
        <Button
          type="submit"
          variant="primary"
          onClick={sendKleos}
          disabled={!buttonEnabled}
        >
          {isGiven ? `Send` : `Give Kudos`}
        </Button>
      </div>
    </div>
  );
};

export default GiveKleosForm;
