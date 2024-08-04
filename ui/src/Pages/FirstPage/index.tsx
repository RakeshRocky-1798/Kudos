import { FC, useState } from 'react';

import WelcomeFirstPage from '@src/Pages/FirstPage/WelcomeFirstPage';
import { getFromLocalStorage } from '@src/service/storage';
import styles from '@src/Pages/FirstPage/WelcomeFirstPage.module.scss';

import FirstGiveKudos from './FirstGiveKudos';

interface FirstPageProps {
  userId: string;
}
const userName: string = getFromLocalStorage('user-name') || '';

const FirstPage: FC<FirstPageProps> = ({ userId }) => {
  const [displayInitForm, setDisplayInitForm] = useState<boolean>(false);

  const handleDisplayInitForm = (display: boolean): void => {
    setDisplayInitForm(display);
  };

  if (displayInitForm) {
    return <FirstGiveKudos userId={userId} />;
  }

  return (
    <div className={styles.wrapper}>
      <div className={styles.container}>
        <WelcomeFirstPage updateDisplayInitForm={handleDisplayInitForm} />
      </div>
    </div>
  );
};

export default FirstPage;
