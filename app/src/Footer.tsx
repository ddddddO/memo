import * as React from 'react';
import Paper from '@mui/material/Paper';
import BottomNavigation from '@mui/material/BottomNavigation';
import BottomNavigationAction from '@mui/material/BottomNavigationAction';
import LockOutlinedIcon from '@mui/icons-material/LockOutlined';
import StickyNote2OutlinedIcon from '@mui/icons-material/StickyNote2Outlined';
import LocalOfferOutlinedIcon from '@mui/icons-material/LocalOfferOutlined';
import { Link } from 'react-router-dom';

export default class Footer extends React.Component {
    render() {
        return (
        <Paper sx={{ position: 'fixed', bottom: 0, left: 0, right: 0 }} elevation={3}>
        <BottomNavigation
          showLabels
        //   value={value}
        //   onChange={(event, newValue) => {
        //     setValue(newValue);
        //   }}
        >
          <BottomNavigationAction
            label="SignIn"
            icon={<LockOutlinedIcon />}
            component={Link}
            to="/signin"
          />
          <BottomNavigationAction
            label="Memos"
            icon={<StickyNote2OutlinedIcon />}
            component={Link}
            to="/memos"
          />
          <BottomNavigationAction
            label="Tags"
            icon={<LocalOfferOutlinedIcon />}
            component={Link}
            to="/memos"
          />
        </BottomNavigation>
      </Paper>
        )
    };
}