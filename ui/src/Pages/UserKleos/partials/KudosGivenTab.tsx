import { FC } from 'react';
import InfiniteScroll from 'react-infinite-scroller';
import { useSelector } from 'react-redux';

import { Typography } from '@navi/web-ui/lib/primitives';
import FallbackComponent from '@src/components/Fallback';
import { selectGivenTabData } from '@src/store/UserKudos/UserSlice';

import useUserKudosApi from '../useUserKudosApi';
import KudosData from './KudosData';

const KudosGivenTab: FC = () => {
  const givenTabData = useSelector(selectGivenTabData);
  const hasMore = givenTabData?.pages?.hasData;
  const { fetchUserKudosData } = useUserKudosApi();

  const handleFetchMore = (): void => {
    fetchUserKudosData(givenTabData?.pages?.pageNumber + 1, 'given', true);
  };

  return (
    <InfiniteScroll
      loadMore={handleFetchMore}
      hasMore={hasMore}
      loader={<FallbackComponent listLoader={true} />}
      useWindow={false}
    >
      <KudosData />
      {!hasMore && givenTabData?.pages?.pageNumber > 0 ? (
        <Typography
          variant="h5"
          color="var(--navi-color-gray-c3)"
          style={{ textAlign: 'center' }}
        >
          You’ve reached the end of the list!
        </Typography>
      ) : null}
    </InfiniteScroll>
  );
};

export default KudosGivenTab;
