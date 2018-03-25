import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import FlatButton from 'material-ui/FlatButton';
import { showNotification as showNotificationAction } from 'admin-on-rest';
import { push as pushAction } from 'react-router-redux';
import request from 'superagent';
import { MyConfig } from './MyConfig';


class ActionExtractButton extends Component {
    handleClick = () => {
        const token = localStorage.getItem('token');
        const { push, showNotification } = this.props;
        request
            .put(MyConfig.API_URL + '/admin/api/v1/actionextract') 
            .set('Authorization', `Bearer ${token}`)
            .then( (response)  => { 
                showNotification('Extracts done'); 
                /*console.log(response); */
                push('/ips'); 
            })
            .catch((e) => {
                // console.error(e);
                showNotification('Error in action extract ', 'warning');
            });

    }

    render() {
        return <FlatButton primary label="Extract" onClick={this.handleClick} />;
    }
}

ActionExtractButton.propTypes = {
    push: PropTypes.func,
    showNotification: PropTypes.func,
};


export default connect(null, {
    showNotification: showNotificationAction,
    push: pushAction,
})(ActionExtractButton);

