/* eslint-disable */
// TODO: Fix is this file without any TS eslint errors.
import { FC } from 'react';
import { useSelector } from 'react-redux';
import { VictoryLabel, VictoryPie } from 'victory';
import emojione from 'emojione';

import { Typography } from '@navi/web-ui/lib/primitives';
import { selectDashboardData } from '@src/store/Dashboard/DashboardSlice';
import { AchievementData } from '@src/Pages/KleosDashboard/types';
import styles from './TotalAchievements.module.scss';

const TotalAchievements: FC = () => {
  const myKleosData: AchievementData[] =
    useSelector(selectDashboardData).totalAchievement || [];
  if (!myKleosData.length) {
    return <Typography variant="p3">No data found</Typography>;
  }
  const totalAchievement = myKleosData.map((item: AchievementData) => ({
    x: item['achievementName'],
    y: item['count'],
  }));

  const emoji = myKleosData.map((item: any) => {
    if (item['emoji'] === 'sports_medal') {
      return 'medal';
    }
    return item['emoji'];
  });

  const colors = [
    'var(--navi-color-yellow-light)',
    'var(--navi-color-blue-light)',
    'var(--navi-color-green-light)',
    '#FF8042',
  ];
  const totalCount = totalAchievement.reduce(
    (total: any, item: { y: any }) => total + item.y,
    0,
  );

  return (
    <>
      <div className={styles['typoWrapper']}>
        <Typography
          variant={'p4'}
          className={styles['totalAchievementsTypo']}
          color={'#969696'}
        >
          Appreciation Roundup
        </Typography>
      </div>

      <div className={styles['pie-container']}>
        <div className={styles['pieChartWrapper']}>
          <svg viewBox="0 0 400 400">
            <VictoryPie
              standalone={false}
              width={400}
              height={400}
              data={totalAchievement}
              colorScale={colors}
              innerRadius={80}
              labelRadius={100}
              labels={() => null}
              padAngle={1}
            />
            <VictoryLabel
              textAnchor="middle"
              style={{ fontSize: 30 }}
              x={200}
              y={200}
              text={totalCount}
            />
          </svg>
        </div>
      </div>

      <div className={styles['labelWrapper']}>
        <div className={styles['legends-container']}>
          <div className={styles['legend-top']}>
            {totalAchievement.length &&
              totalAchievement
                .slice(0, 2)
                .map((achievement: any, index: any) => (
                  <div className={styles['legend-item']} key={index}>
                    <div className={styles['legWrapper']}>
                      <div className={styles['legend-color2']}>
                        <div
                          className={styles['legend-color2']}
                          style={{ backgroundColor: colors[index] }}
                        ></div>
                      </div>
                      <div className={styles['legend-value']}>
                        {achievement['y']}
                      </div>
                    </div>
                    <div className={styles['legend-label']}>
                      {achievement['x']}{' '}
                      {emojione.shortnameToUnicode(`:${emoji[index]}:`)}
                    </div>
                  </div>
                ))}
          </div>

          <div className={styles['legend-bottom']}>
            {totalAchievement.length &&
              totalAchievement.slice(2).map((achievement: any, index: any) => (
                <div className={styles['legend-item']} key={index}>
                  <div className={styles['legWrapper']}>
                    <div className={styles['legend-color2']}>
                      <div
                        className={styles['legend-color2']}
                        style={{ backgroundColor: colors[index + 2] }}
                      ></div>
                    </div>
                    <div className={styles['legend-value']}>
                      {achievement['y']}
                    </div>
                  </div>
                  <div className={styles['legend-label']}>
                    {achievement['x']}{' '}
                    {emojione.shortnameToUnicode(`:${emoji[index + 2]}:`)}
                  </div>
                </div>
              ))}
          </div>
        </div>
      </div>
    </>
  );
};

export default TotalAchievements;
