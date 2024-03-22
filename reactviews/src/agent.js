import React from 'react';
import { List, Datagrid, TextField, Edit, Create, SimpleForm,
    TextInput, required, EditButton, DateField,
    RichTextField, SelectInput, Filter, Labeled,
    SimpleList, useRecordContext
} from 'react-admin';
import { RichTextInput } from 'ra-input-rich-text';
import { useMediaQuery } from '@mui/material';
import { roles } from './MyConfig';


const cmd = [
    { label: 'Send Lines', id: 'SendLines', name: 'SendLines' },
    { label: 'Search Full File', id: 'FullSearch', name: 'FullSearch' },
    { label: 'STOP', id: 'STOP', name: 'STOP' },
];

const ColoredTextField = (props) => {
    const record = useRecordContext();
    return (
        <TextField
            sx={{ color: record.status === 'OnLine' ? 'green' : 'red' }}
            {...props}
        />
    );
};
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

export const AgentList = (props) => {
    const isSmall = useMediaQuery(
        theme => theme.breakpoints.down('sm'),
        { noSsr: true }
    );
    return (
        <List bulkActionButtons={false} filters={<AgentFilter />} {...props}>
            {isSmall ? (
                <SimpleList
                    primaryText={record => `[${record.status}] ${record.ip}`}
                    secondaryText={record => record.filesurvey}
                    tertiaryText={record => new Date(record.updated).toLocaleString()}
                />
            ) : (
                <Datagrid>
                    <TextField source="ip" />
                    <TextField source="filesurvey" />
                    <ColoredTextField source="status" />
                    <TextField source="role" />
                    <RichTextField source="lines" sx = { styles.field }/>
                    <RichTextField source="comment" stripTags />
                    <TextField source="cmd" />
                    <DateField label="updated" source="updated" showTime />
                    <EditButton />
                </Datagrid>
            )
            }
        </List>
    );
};


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
            <Labeled label="File Survey">
                <TextField source="filesurvey" />
            </Labeled>
            <Labeled label="Created">
                <DateField label="created" source="created" showTime />
            </Labeled>
            <Labeled label="Updated">
                <DateField label="updated" source="updated" showTime />
            </Labeled>
            <Labeled label="Lines">
                <TextField source="lines"  style={{width: '100%', whiteSpace: 'pre-line'}}/>
            </Labeled>
            <SelectInput source="role" choices={roles} />
            <SelectInput source="cmd" choices={cmd} optionText="label" />
            <RichTextInput source="comment" />
            <TextInput source="crca" disabled />
        </SimpleForm>
    </Edit>
);
