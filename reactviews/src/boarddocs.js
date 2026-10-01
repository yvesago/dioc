import React from 'react';
import { Edit, SimpleForm, Toolbar, SaveButton  } from 'react-admin';
//import { RichTextInput } from 'ra-input-rich-text';
const RichTextInput = React.lazy(() =>
    import('ra-input-rich-text').then(module => ({
        default: module.RichTextInput,
    }))
);
import Divider from '@mui/material/Divider';


const BoardEditToolbar = () => {
    return (
        <Toolbar>
            <SaveButton />
        </Toolbar>
    );
};

export const BoardEdit = () => (
    <Edit title="Docs" redirect="/" undoable={false}>
        <SimpleForm toolbar={<BoardEditToolbar />}>
            <Divider flexItem />
            <RichTextInput
                source="docs"
                label="Main Doc"
            />
            <br />
            <br />
            <Divider flexItem />
            <RichTextInput
                source="docagents"
                label="Agents"
            />
            <br />
            <br />
            <Divider flexItem />
            <RichTextInput
                source="docsearchs"
                label="Searchs"
            />
        </SimpleForm>
    </Edit>
);
