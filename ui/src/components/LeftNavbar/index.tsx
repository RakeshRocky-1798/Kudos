import React, { useState } from 'react';
import { useAuth } from '@src/hooks/useAuth';
import { USER_PERMISSION } from '@src/Pages/KleosDashboard/constants';

import { DarkKnightTokenParsed } from '@navi/dark-knight';
import { Navbar } from '@navi/web-ui/lib/components';
import NaviNewLogoIcon from '@navi/web-ui/lib/icons/NaviLogoIcon/NaviNewLogoIcon';
import DashboardIcon from '@navi/web-ui/lib/icons/DashboardIcon';
import { Typography, Avatar } from '@navi/web-ui/lib/primitives';
import styles from './LeftNav.module.scss';
import { NavItemType } from '@navi/web-ui/lib/components/Navbar/types';
import cx from 'classnames';
import '@components/LeftNavbar/leftnav.css';
import { useNavigate } from 'react-router-dom';
import AdminPortalIcon from '@src/assets/icons/videocam.svg';
import StarIcon from '@src/assets/icons/star.svg';
import InfoIcon from '@src/assets/icons/infoIcon.svg';

interface LeftNavProps {
  children?: React.ReactNode;
}

const LeftNav: React.FC<LeftNavProps> = ({ children }) => {
  const navigate = useNavigate();
  const [isExpanded, setIsExpanded] = useState(false);
  const handleExpandChange = (isExpanded: boolean): void => {
    setIsExpanded(isExpanded);
  };

  const navItemList: NavItemType[] = [
    {
      itemType: 'simpleNavItem',
      label: 'Dashboard',
      route: '/',
      Icon: () => <DashboardIcon width={24} height={24} color={'#FCFCFD'} />,
      handleNavigation: () => navigate('/'),
    },
    {
      itemType: 'simpleNavItem',
      label: 'User Kudos',
      route: '/user-kudos',
      Icon: () => <img src={StarIcon} alt="giveKleosAlt" />,
      handleNavigation: () => navigate('/user-kudos'),
    },
  ];

  const { hasPermission } = useAuth();
  // const bottomNavItemList: bottomNavItemType[] = [];
  const bottomNavItemList: NavItemType[] = [];

  bottomNavItemList.push({
    itemType: 'simpleNavItem',
    label: 'About Kudos',
    route: '/about-kudos',
    Icon: () => (
      <img
        src={InfoIcon}
        alt="info icon"
        style={{ width: '24px', height: '24px' }}
      />
    ),
    handleNavigation: () => navigate('/about-kudos'),
  });

  if (hasPermission(USER_PERMISSION.REQUEST_ADMIN_BOARD_READ_PERMISSION)) {
    bottomNavItemList.push({
      itemType: 'simpleNavItem',
      label: 'Admin',
      route: '/admin',
      Icon: () => <img src={AdminPortalIcon} alt="giveKleosAlt" />,
      handleNavigation: () => navigate('/admin'),
    });
  }

  // const bottomNavItemList: bottomNavItemType[] = [
  //   {
  //     itemType: 'simpleNavItem',
  //     label: 'Admin',
  //     route: '/user-kudos',
  //     Icon: () => <img src={AdminPortalIcon} alt="giveKleosAlt" />,
  //     handleNavigation: () => navigate('/admin'),
  //   },
  // ];

  const [dialogOpen, setDialogOpen] = useState(false);

  const closeDialog = (): void => setDialogOpen(false);

  const confirmLogout = (): void => {
    closeDialog();
    localStorage.removeItem('auth-token');
    localStorage.removeItem('user-data');
    window.location.reload();
  };

  const returnProfileImage = (): JSX.Element | undefined => {
    const userData: string | null = localStorage.getItem('user-data');
    if (userData) {
      const parsedUserData: DarkKnightTokenParsed = JSON.parse(userData);
      return (
        <Avatar
          isImage
          src={parsedUserData.profilePictureUrl}
          alt="profileImage"
          sizeVariant="md"
        />
      );
    }
    return undefined;
  };

  const footerDetails = {
    footerText: localStorage.getItem('user-name') || 'User',
    avatar: returnProfileImage(),
    options: [
      {
        label: 'Logout',
        handleClick: confirmLogout,
      },
    ],
  };
  const headerContent: JSX.Element = (
    <div>
      <div className={styles.header}>
        <NaviNewLogoIcon width={28} height={28} />
        {isExpanded && (
          <Typography variant="h3" color="white" style={{ marginTop: '8px' }}>
            Kudos UI
          </Typography>
        )}
      </div>
    </div>
  );

  return (
    <>
      <Navbar
        headerContent={headerContent}
        navItemList={navItemList}
        bottomNavItemList={bottomNavItemList}
        footer={footerDetails}
        handleNavbarStateChange={handleExpandChange}
      >
        <div
          className={cx(
            styles['page-content-container'],
            isExpanded ? styles['expanded-page-content'] : '',
          )}
        >
          {children}
        </div>
      </Navbar>
    </>
  );
};

export default LeftNav;
