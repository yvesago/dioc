import React from 'react';
import CardActions from '@material-ui/core/CardActions';
import { List, Datagrid, TextField, Edit, Create, SimpleForm,
    TextInput, required, EditButton, DateField, CreateButton,
    RichTextField, SelectInput, Filter, Responsive, SimpleList,
    BooleanInput, BooleanField, DateInput, RefreshButton
} from 'react-admin';
import RichTextInput from 'ra-input-rich-text';
import { withStyles } from '@material-ui/core/styles';


import ActionExtractButton from './ActionExtract';
import { roles } from './MyConfig';

const actions = [
    { name: 'AddIP', id: 'AddIP' },
    { name: 'Compress', id: 'Compress' },
    { name: 'Delete', id: 'Delete' },
];

const colored = WrappedComponent => props => props.record.active === true ?
    <span style={{ color: 'green', fontWeight: 'bold' }}><WrappedComponent {...props} /></span> :
    <WrappedComponent {...props} />;


const cardActionStyle = {
    zIndex: 2,
    display: 'inline-block',
    float: 'right',
};

const ExtractActions = ({ resource, filters, displayedFilters, filterValues, basePath, showFilter }) => (
    <CardActions style={cardActionStyle}>
        {filters && React.cloneElement(filters, { resource, showFilter, displayedFilters, filterValues, context: 'button' }) }
        <ActionExtractButton />
        <CreateButton basePath={basePath} />
        <RefreshButton />
    </CardActions>
);


const ColoredTextField = colored(TextField);


const ExtractFilter = (props) => (
    <Filter {...props}>
        <TextInput label="Comment" source="comment" />
        <SelectInput source="role" choices={roles} allowEmpty alwaysOn />
    </Filter>
);

const styles = {
    field: {
        width: '250px', whiteSpace: 'nowrap', overflow: 'hidden', textOverflow: 'ellipsis'}
};

export const ExtractList = withStyles(styles)(({ classes, ...props }) => (
    <List bulkActions={false} filters={<ExtractFilter />} perPage={30} actions={<ExtractActions />} {...props}>
        <Responsive
            small={
                <SimpleList
                    primaryText={record => record.search}
                    secondaryText={record => `Role: ${record.role}, Active: ${record.active}`}
                    tertiaryText={record => new Date(record.updated).toLocaleString()}
                />
            }
            medium={
                <Datagrid>
                    <TextField source="search" />
                    <TextField source="role" />
                    <ColoredTextField source="action" />
                    <BooleanField source="active" />
                    <RichTextField source="comment" className={classes.field} stripTags />
                    <DateField label="updated" source="updated" showTime />
                    <EditButton />
                </Datagrid>
            }
        />
    </List>
));


export const ExtractCreate = (props) => (
    <Create {...props}>
        <SimpleForm>
            <TextInput source="search" validate={required()} />
            <SelectInput source="role" choices={roles} allowEmpty />
            <DateInput label="From" source="fromdate" parse={dateParser} allowEmpty />
            <DateInput label="To" source="todate" parse={dateParser} allowEmpty />
            <SelectInput source="action" choices={actions} allowEmpty />
            <BooleanInput source="active" />
            <RichTextInput source="comment" />
        </SimpleForm>
    </Create>
);


export const ExtractEdit = (props) => (
    <Edit  {...props}>
        <SimpleForm>
            <TextInput source="search" validate={required()} />
            <SelectInput source="role" choices={roles} allowEmpty />
            <DateInput label="From" source="fromdate" parse={dateParser} allowEmpty />
            <DateInput label="To" source="todate" parse={dateParser} allowEmpty />
            <SelectInput source="action" choices={actions} allowEmpty />
            <BooleanInput source="active" />
            <RichTextInput source="comment" />
            <DateField label="Created" source="created" showTime />
            <DateField label="Updated" source="updated" showTime />
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
