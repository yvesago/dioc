import React from 'react';
import { List, Datagrid, TextField, Edit, Create, SimpleForm,
    TextInput, required, EditButton, DateField,
    RichTextField, SelectInput, Filter, Responsive,
    SimpleList
} from 'react-admin';
import RichTextInput from 'ra-input-rich-text';
import { withStyles } from '@material-ui/core/styles';
import classnames from 'classnames';
import { roles } from './MyConfig';


const cmd = [
    { label: 'Send Lines', id: 'SendLines', name: 'SendLines' },
    { label: 'Search Full File', id: 'FullSearch', name: 'FullSearch' },
    { label: 'STOP', id: 'STOP', name: 'STOP' },
];

const coloredStyles = {
    up: { color: 'green' },
    down: { color: 'red' },
};

const ColoredTextField = withStyles(coloredStyles)(
    ({ classes, ...props }) => (
        <TextField
            className={classnames({
                [classes.down]: props.record.status === 'OffLine',
                [classes.up]: props.record.status === 'OnLine',
            })}
            {...props}
        />
    ));

ColoredTextField.defaultProps = TextField.defaultProps;


const AgentTitle = ({ record }) => {
    return <span>Agent {record ? `${record.ip}` : ''}</span>;
};


const AgentFilter = (props) => (
    <Filter {...props}>
        <TextInput label="Ip" source="ip" alwaysOn />
        <TextInput label="Lines" source="lines" />
        <TextInput label="Comment" source="comment" />
        <SelectInput source="role" choices={roles} />
    </Filter>
);

const styles = {
    field: {
        display: 'inline-block', width: '200px', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis'}
};

export const AgentList = withStyles(styles)(({ classes, ...props }) => (
    <List bulkActionButtons={false} filters={<AgentFilter />} {...props}>
        <Responsive
            small={
                <SimpleList
                    primaryText={record => `[${record.status}] ${record.ip}`}
                    secondaryText={record => record.filesurvey}
                    tertiaryText={record => new Date(record.updated).toLocaleString()}
                />
            }
            medium={
                <Datagrid>
                    <TextField source="ip" />
                    <TextField source="filesurvey" />
                    <ColoredTextField source="status" />
                    <TextField source="role" />
                    <RichTextField source="lines" className={classes.field} />
                    <RichTextField source="comment" stripTags />
                    <TextField source="cmd" />
                    <DateField label="updated" source="updated" showTime />
                    <EditButton />
                </Datagrid>
            }
        />
    </List>
));


export const AgentCreate = (props) => (
    <Create {...props}>
        <SimpleForm redirect="list">
            <TextInput source="ip" validate={required()} />
            <TextInput source="filesurvey" />
            <TextInput source="status" />
            <SelectInput source="role" choices={roles} allowEmpty />
            <TextInput source="lines" options={{ multiLine: true }}  />
            <TextInput multiline source="comment" />
        </SimpleForm>
    </Create>
);

export const AgentEdit = (props) => (
    <Edit title={<AgentTitle />} {...props}>
        <SimpleForm>
            <TextField source="filesurvey" />
            <DateField label="created" source="created" showTime />
            <DateField label="updated" source="updated" showTime />
            <TextField source="lines"  style={{width: '100%', whiteSpace: 'pre-line'}}/>
            <SelectInput source="role" choices={roles} allowEmpty />
            <SelectInput source="cmd" choices={cmd} allowEmpty optionText="label" />
            <RichTextInput source="comment" />
            <TextInput source="crca" disabled />
        </SimpleForm>
    </Edit>
);
