import { FC, useState, ReactNode } from 'react';
import Lottie from 'react-lottie';
import { useSelector } from 'react-redux';
import { useNavigate } from 'react-router-dom';

import Button from '@navi/web-ui/lib/primitives/Button';
import Typography from '@navi/web-ui/lib/primitives/Typography';
import { Grid } from '@navi/web-ui/lib/layouts/index';

import OnboardingAnimation1 from '@src/assets/lotties/OnboardingAnimation_1.json';
import OnboardingAnimation2 from '@src/assets/lotties/OnboardingAnimation_2.json';

import ErrorBoundary from '@src/components/ErrorBoundary/ErrorBoundary';
import FallbackComponent from '@src/components/Fallback';
import { selectIsGiven } from '@src/store/Dashboard/DashboardSlice';
import styles from '@src/Pages/FirstPage/WelcomeFirstPage.module.scss';
import { getFromLocalStorage } from '@src/service/storage';

/* Extract Grid components */
const { GridContainer, GridRow, GridColumn } = Grid;

interface WelcomeFirstPageProps {
  updateDisplayInitForm?: (display: boolean) => void;
}

interface CarouselContentProps {
  title: string;
  subTitle: string;
  descriptions: string[];
  linkDetails?: ReactNode;
}

const userName: string = getFromLocalStorage('user-name') || '';

const returnGuideLineLinkComponent = (): ReactNode => {
  return (
    <Typography
      variant="p2"
      className={styles['content-item']}
      color={'var(--navi-color-navigation-blue-border)'}
    >
      Click{' '}
      <a
        href="https://drive.google.com/drive/folders/10A7weM_B7ArKKEL3piwF6JBTp6iOdgEG?usp=sharing"
        target="_blank"
        rel="noreferrer"
      >
        here
      </a>{' '}
      to know more about Kudos.
    </Typography>
  );
};

const carouselContentPartOne: CarouselContentProps = {
  title: `Hey, ${userName} 👋`,
  subTitle: 'Welcome to kudos',
  descriptions: [
    `Navi's very own peer-to-peer appreciation tool designed to foster a culture of recognition and celebrate excellence within our team.`,
    `Get ready to hop on your Kudos journey and give Kudos to your rockstar shipmates!`,
  ],
  linkDetails: returnGuideLineLinkComponent(),
};

const carouselContentPartTwo: CarouselContentProps = {
  title: 'Just a friendly heads-up:',
  subTitle: '',
  descriptions: [
    `Our new Kudos tool is all about spreading positivity and giving virtual high fives to each other for a job well done!`,
    `Remember, these kudos are purely for recognition and won't impact individual performance reviews. Let's keep the appreciation flowing freely!`,
  ],
  linkDetails: returnGuideLineLinkComponent(),
};

const WelcomeFirstPage: FC<WelcomeFirstPageProps> = ({
  updateDisplayInitForm,
}) => {
  const [carouselIndex, setCarouselIndex] = useState<number>(0);
  const isGiven = useSelector(selectIsGiven);
  const navigate = useNavigate();

  const returnLottieContent = (index: number): JSX.Element => {
    const defaultOptions = {
      loop: true,
      autoplay: true,
      animationData: index === 0 ? OnboardingAnimation1 : OnboardingAnimation2,
      rendererSettings: {
        scaleMode: 'fit',
        preserveAspectRatio: 'xMidYMid slice',
      },
    };

    if (index === 1) {
      return (
        <div className={styles['lottie-with-background']}>
          <Lottie
            options={defaultOptions}
            style={{ width: 'inherit', height: 'inherit' }}
          />
        </div>
      );
    }

    return (
      <div className={styles['lottie']}>
        <Lottie
          options={defaultOptions}
          style={{ width: 'inherit', height: 'inherit' }}
        />
      </div>
    );
  };

  const handleNavigatorClick = (goBack = false): void => {
    if (goBack) {
      setCarouselIndex(0);
      return;
    }

    if (carouselIndex === 0) {
      setCarouselIndex(1);
      return;
    }

    if (carouselIndex === 1) {
      if (updateDisplayInitForm) updateDisplayInitForm(true);
      if (isGiven) navigate('/');
      return;
    }
  };

  const carouselContent = (index: number): JSX.Element => {
    const content =
      index === 0 ? carouselContentPartOne : carouselContentPartTwo;

    return (
      <GridContainer>
        <GridRow>
          <GridColumn sm={6} className={styles['content-banner']}>
            {returnLottieContent(index)}
          </GridColumn>
          <GridColumn sm={6}>
            <Typography
              variant="h1"
              className={styles['content-item']}
              color={'var(--navi-color-gray-c1)'}
            >
              {content.title}
            </Typography>
            {content.subTitle.length > 0 && (
              <Typography
                variant="h2"
                className={styles['content-item']}
                color={'var(--navi-color-gray-c1)'}
              >
                Welcome to kudos
              </Typography>
            )}
            {content.descriptions.map((description, index) => (
              <Typography
                key={`description-${index}`}
                variant="p2"
                className={styles['content-item']}
                color={'var(--navi-color-navigation-blue-border)'}
              >
                {description}
              </Typography>
            ))}
            {content?.linkDetails !== null && content?.linkDetails}
          </GridColumn>
        </GridRow>
        <GridRow>
          <GridColumn sm={6}>{''}</GridColumn>
          <GridColumn sm={6}>
            {carouselIndex === 1 && (
              <Button
                variant="secondary"
                onClick={(): void => handleNavigatorClick(true)}
                className={styles['back-button']}
              >
                Back
              </Button>
            )}
            <Button
              type="submit"
              variant="primary"
              className={styles['navigate-button']}
              onClick={(): void => handleNavigatorClick()}
            >
              {carouselIndex === 0 ? 'Next' : 'Continue'}
            </Button>
          </GridColumn>
        </GridRow>
      </GridContainer>
    );
  };

  const returnContent = (): JSX.Element => {
    return <section>{carouselContent(carouselIndex)}</section>;
  };

  return (
    <ErrorBoundary fallbackComponent={<FallbackComponent />}>
      {returnContent()}
    </ErrorBoundary>
  );
};

export default WelcomeFirstPage;
