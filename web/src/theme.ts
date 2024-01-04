import type { ThemeConfig } from 'antd';

const theme: ThemeConfig = {
  cssVar: false,
  hashed: true,
  token: {
    colorPrimary: '#13C2C2',
    borderRadiusLG: 4,
  },
  components: {
    Layout: {
      headerHeight: 56,
      headerBg: 'rgba(255, 255, 255, 0.6)',
      headerPadding: '0 12px',
    },
    Menu: {
      itemBg: 'transparent',
      activeBarBorderWidth: 0,
      // activeBarHeight: 0,
    },
    Tabs: {
      horizontalItemGutter: 24,
    }
  },
};

export default theme;