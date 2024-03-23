import React from 'react';
import { useMediaQuery } from '@mui/material';
import { List, Datagrid, TextField, Edit, Create, SimpleForm,
    TextInput, required, EditButton, DateField, CreateButton,
    RichTextField, SelectInput, Filter, SimpleList, useRecordContext,
    BooleanInput, BooleanField, DateInput, RefreshButton,
    TopToolbar, FilterButton, Labeled
} from 'react-admin';
//import { RichTextInput } from 'ra-input-rich-text';
const RichTextInput = React.lazy(() =>
    import('ra-input-rich-text').then(module => ({
        default: module.RichTextInput,
    }))
);


import ActionExtractButton from './ActionExtract';
import { roles } from './MyConfig';

const actions = [
    { name: 'AddIP', id: 'AddIP' },
    { name: 'Compress', id: 'Compress' },
    { name: 'Delete', id: 'Delete' },
];


const ColoredTextField = (props) => {
    const record = useRecordContext();
    return (
        <TextField
            sx={ record.active === true ? { color: 'green', fontWeight: 'bold' } : {} }
            {...props}
        />
    );
};
ColoredTextField.defaultProps = TextField.defaultProps;


const ExtractActions = () => (
    <TopToolbar>
        <FilterButton />
        <ActionExtractButton />
        <CreateButton />
    </TopToolbar>
);

const postFilters = [
    <TextInput label="Comment" source="comment" />,
    <SelectInput source="role" choices={roles} alwaysOn />,
];

const styles = {
    field: {
        display: 'inline-block', width: '250px', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis'}
};

export const ExtractList = ({ props }) => {
    const isSmall = useMediaQuery(
        theme => theme.breakpoints.down('sm'),
        { noSsr: true }
    );
    return (
        <List bulkActionButtons={false} filters={postFilters} perPage={50} actions={<ExtractActions />} sort={{ field: 'updated', order: 'DESC' }} {...props}>
            {isSmall ? (
                <SimpleList
                    primaryText={record => record.search}
                    secondaryText={record => `Role: ${record.role}, Active: ${record.active}`}
                    tertiaryText={record => new Date(record.updated).toLocaleString()}
                />
            ) : (
                <Datagrid>
                    <TextField source="search" />
                    <TextField source="role" />
                    <ColoredTextField source="action" />
                    <BooleanField source="active" />
                    <RichTextField source="comment" sx={styles.field} stripTags />
                    <DateField label="updated" source="updated" showTime />
                    <EditButton />
                </Datagrid>
            )}
        </List>
    );
};


export const ExtractCreate = (props) => (
    <Create {...props}>
        <SimpleForm redirect="list">
            <TextInput source="search" validate={required()} />
            <SelectInput source="role" choices={roles} />
            <DateInput label="From" source="fromdate" parse={dateParser} />
            <DateInput label="To" source="todate" parse={dateParser} />
            <SelectInput source="action" choices={actions} />
            <BooleanInput source="active" />
            <RichTextInput source="comment" />
        </SimpleForm>
    </Create>
);


export const ExtractEdit = (props) => (
    <Edit  {...props}>
        <SimpleForm>
            <TextInput source="search" validate={required()} />
            <SelectInput source="role" choices={roles} />
            <DateInput label="From" source="fromdate" parse={dateParser} />
            <DateInput label="To" source="todate" parse={dateParser} />
            <SelectInput source="action" choices={actions} />
            <BooleanInput source="active" />
            <RichTextInput source="comment" />
            <Labeled label="Created">
                <DateField source="created" showTime />
            </Labeled>
            <Labeled label="Updated">
                <DateField source="updated" showTime />
            </Labeled>
        </SimpleForm>
    </Edit>
);

/*const dateFormatter = v => { // from record to input
    // v is a `Date` object
    if (!(v instanceof Date) || isNaN(v)) return;
    const pad = '00';
    const yy = v.getFullYear().toString();
    const mm = (v.getMonth() + 1).toString();
    const dd = v.getDate().toString();
    if (yy === "0001") return "";
    return `${yy}-${(pad + mm).slice(-2)}-${(pad + dd).slice(-2)}`;
};*/

const dateParser = v => { // from input to record
    // v is a string of "YYYY-MM-DD" format
    const match = /(\d{4})-(\d{2})-(\d{2})/.exec(v);
    if (match === null) return;
    const d = new Date(match[1], parseInt(match[2], 10) - 1, match[3]);
    if (isNaN(d)) return;
    return d;
};
