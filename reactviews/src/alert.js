import React from 'react';
import { List, Datagrid, TextField, DateField, EditButton,
    Edit, SimpleForm, Filter, TextInput, SelectInput,
    RichTextField, Responsive, SimpleList
} from 'react-admin';
import RichTextInput from 'ra-input-rich-text';
import { withStyles } from '@material-ui/core/styles';
import classnames from 'classnames';
import { roles } from './MyConfig';

const coloredStyles = {
    warn: { color: 'orange' },
    critic: { color: 'red' },
};

const ColoredTextField = withStyles(coloredStyles)(
    ({ classes, ...props }) => (
        <TextField
            className={classnames({
                [classes.warn]: props.record.level === 'warn',
                [classes.critic]: props.record.level === 'critic',
            })}
            {...props}
        />
    ));

ColoredTextField.defaultProps = TextField.defaultProps;


const AlertFilter = (props) => (
    <Filter {...props}>
        <TextInput label="Line" source="line" />
        <TextInput label="Comment" source="comment" />
        <SelectInput source="role" choices={roles} allowEmpty alwaysOn />
    </Filter>
);

const styles = {
    field: {
        display: 'inline-block', width: '250px', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis'}
};

export const AlertList = withStyles(styles)(({ classes, ...props }) => (
    <List filters={<AlertFilter />}  sort={{ field: 'updated', order: 'DESC' }} perPage={30} {...props}>
        <Responsive
            small={
                <SimpleList
                    primaryText={record => `[${record.level}] ${record.search}`}
                    secondaryText={record => record.line}
                    tertiaryText={record => new Date(record.updated).toLocaleString()}
                />
            }
            medium={
                <Datagrid>
                    <TextField source="search" />
                    <TextField source="role" />
                    <ColoredTextField source="level" />
                    <TextField source="line" />
                    <RichTextField source="comment" className={classes.field} stripTags />
                    <DateField label="updated" source="updated" showTime />
                    <EditButton />
                </Datagrid>
            }
        />
    </List>
));

export const AlertEdit = (props) => (
    <Edit {...props}>
        <SimpleForm>
            <TextField source="search" />
            <TextField source="ip" label="Agent" />
            <TextField source="filesurvey" />
            <TextField source="role" />
            <ColoredTextField source="level" />
            <TextField source="line" style={{width: '100%', whiteSpace: 'pre-line'}}/>
            <RichTextInput source="comment" />
            <DateField label="created" source="created" showTime />
            <DateField label="updated" source="updated" showTime />
        </SimpleForm>
    </Edit>
);
