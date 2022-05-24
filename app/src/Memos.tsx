import * as React from 'react';
import Box from '@mui/material/Box';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import ListItemAvatar from '@mui/material/ListItemAvatar';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import Avatar from '@mui/material/Avatar';
import IconButton from '@mui/material/IconButton';
import FormGroup from '@mui/material/FormGroup';
import FormControlLabel from '@mui/material/FormControlLabel';
import Checkbox from '@mui/material/Checkbox';
import Grid from '@mui/material/Grid';
import Typography from '@mui/material/Typography';
import Container from '@mui/material/Container';
import FolderIcon from '@mui/icons-material/Folder';
import DeleteIcon from '@mui/icons-material/Delete';
import { FixedSizeList, ListChildComponentProps } from 'react-window';
import ListItemButton from '@mui/material/ListItemButton';

export default class Memos extends React.Component {
  renderRow(props: ListChildComponentProps) {
    const { index, style } = props;

    return (
      <ListItem
        style={style} key={index} component="div" disablePadding
        secondaryAction={
          <IconButton edge="end" aria-label="delete">
            <DeleteIcon />
          </IconButton>
        }
      >
        <ListItemAvatar>
          <Avatar>
            <FolderIcon />
          </Avatar>
        </ListItemAvatar>
        <ListItemText
          primary={`Item ${index + 1}`}
          secondary='Secondary text'
        />
      </ListItem>
    );
  }


  generate(element: React.ReactElement) {
    return [0, 1, 2].map((value) =>
      React.cloneElement(element, {
        key: value,
      })
    );
  }

  render() {
    return (
      <Container maxWidth="sm">
        <Box sx={{ width: '100%', height: 400, maxWidth: 360, bgcolor: 'background.paper' }}>
          <FixedSizeList
            height={400}
            width={360}
            itemSize={46}
            itemCount={1000}
            overscanCount={5}
          >
            {this.renderRow}
          </FixedSizeList>
        </Box>
      </Container>
    );
  }
}