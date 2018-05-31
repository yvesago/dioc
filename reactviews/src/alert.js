import React from 'react';
import { List, Datagrid, TextField, DateField, EditButton,
    Edit, SimpleForm, Filter, TextInput, SelectInput, DeleteButton,
    RichTextField, Responsive, SimpleList
} from 'react-admin';
import RichTextInput from 'ra-input-rich-text';
import { withStyles } from '@material-ui/core/styles';
import { roles } from './MyConfig';

const colored = WrappedComponent => props => props.record.level === 'critic' ?
    <span style={{ color: 'red' }}><WrappedComponent {...props} /></span> :
    props.record.level === 'warn' ?
        <span style={{ color: 'orange' }}><WrappedComponent {...props} /></span> :
        <WrappedComponent {...props} />;

const ColoredTextField = colored(TextField);

const AlertFilter = (props) => (
    <Filter {...props}>
        <TextInput label="Line" source="line" />
        <TextInput label="Comment" source="comment" />
        <SelectInput source="role" choices={roles} allowEmpty alwaysOn />
    </Filter>
);

const styles = {
    field: {
        width: '250px', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis'}
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
                    <RichTextField source="comment" className={styles.field} stripTags />
                    <DateField label="updated" source="updated" showTime />
                    <EditButton />
                    <DeleteButton />
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
            <span style={{transform: 'scale(0.75) translate(0px, 10px)', transformOrigin: 'left top 0px',  color: 'rgba(0, 0, 0, 0.3)'}}>Level</span><ColoredTextField source="level" />
            <TextField source="line"  style={{whiteSpace: 'pre-line'}}/>
            <RichTextInput source="comment" />
            <DateField label="created" source="created" showTime />
            <DateField label="updated" source="updated" showTime />
        </SimpleForm>
    </Edit>
);
