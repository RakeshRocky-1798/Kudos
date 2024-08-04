import { FC } from 'react';
import { Typography, Button } from '@navi/web-ui/lib/primitives';
import { ArrowRightUpIcon } from '@navi/web-ui/lib/icons';
import errorBoundary from '@src/assets/icons/errorBoundary.svg';
import commonStyles from './ErrorIndex.module.scss';

const FallbackComponentUI: FC = () => {
  const handleButtonClick = (): void => {
    window.location.href = '/';
  };

  return (
    <div className={commonStyles.centeredContainer}>
      <div className={commonStyles.centeredContent}>
        <img src={errorBoundary} alt="ErrorBoundaryIcon" />
        <Typography className={commonStyles.errorBound} variant="h3">
          Oops, Something Went Wrong!
        </Typography>
        <Button
          startAdornment={<ArrowRightUpIcon color="white" />}
          onClick={handleButtonClick}
        >
          Go To Home
        </Button>
      </div>
    </div>
  );
};

export default FallbackComponentUI;
