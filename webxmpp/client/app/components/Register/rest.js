import axios from 'axios';

const API = 'http://localhost:8080'

export function jidTaken(jid) {
  return axios.get(`${API}/users/${jid}`)
    .then(function (response) {
      return true;
    })
    .catch(function (error) {
      if (error.response.status == 404) {
        return false;
      } else {
        return true;
      }
    });
}

export function createUser(userData) {
  return axios.post(`${API}/users`, userData)
    .then(function (response) {
      return true;
    })
    .catch(function (error) {
        return false;
    });
}
