import { Button, useNotify } from 'react-admin';
import { useNavigate } from 'react-router-dom';
import { Alert } from '@mui/material';
import request from 'superagent';
import { MyConfig } from './MyConfig';


const ActionExtractButton = () => {
    const navigate = useNavigate();
    const notify = useNotify();
    const handleClick = () => {
        const token = localStorage.getItem('gotoken');
        request
            .put(MyConfig.API_URL + '/admin/api/v1/actionextract') 
            .set('Authorization', `Bearer ${token}`)
            .then( (response)  => { 
                notify(<Alert severity="success">Extracts done</Alert>);
                //console.log(response);
                navigate('/ips');
            })
            .catch((e) => {
                //console.error(e);
                notify(<Alert severity="error">Error in action extract</Alert>);
            });

    };

    return <Button label='Extract' onClick={handleClick} />;
};

export default ActionExtractButton;
