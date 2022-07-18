/*
 * Copyright 2022 Google Inc. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the
 * License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
 * express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

import '../node_modules/bootstrap/dist/css/bootstrap.min.css';
import '../node_modules/firebaseui/dist/firebaseui.css';
import '../public/style.css';

// Import Firebase dependencies.
import firebase from 'firebase/compat/app';
import * as firebaseui from 'firebaseui';
import * as ciap from 'gcip-iap';

/** @return Whether the current browser is Safari. */
function isSafari(): boolean {
    const userAgent = navigator.userAgent.toLowerCase();
    return userAgent.indexOf('safari/') !== -1 &&
        userAgent.indexOf('chrome/') === -1 &&
        userAgent.indexOf('crios/') === -1 &&
        userAgent.indexOf('android/') === -1;
}

// The list of UI configs for each supported tenant.
const tenantsConfig = {
    // Project level IdPs flow. * is a wildcard for all tenants
    '*': {
        displayName: 'My Organization',
        signInOptions: [
            "oidc.auth0"
        ],
        // Do not trigger immediate redirect in Safari without some user
        // interaction.
        immediateFederatedRedirect: !isSafari(),
    },
};

const configs = {};

// This is the config for the APIKey that is passed from IAP. It's used to select the auth domain. Alternativly we could overwrite the selectTenant Method of the Handler. https://github.com/firebase/firebaseui-web/blob/1a8f9e523f76df7a6c5da36b63b4d2a5b34b64ef/javascript/widgets/firebaseuihandler.js#L228
configs["apiKey"] = {
    authDomain: "<gcp_project_id>.firebaseapp.com", // authDomain is <gcp_project_id>.firebaseapp.com
    displayMode: 'optionsFirst',
    tenants: tenantsConfig,
};

// This will handle the underlying handshake for sign-in, sign-out,
// token refresh, safe redirect to callback URL, etc.
const handler = new firebaseui.auth.FirebaseUiHandler('#firebaseui-container', configs);
handler.startSignIn = function (auth, selectedTenantInfo) {
    return new Promise((resolve, reject) => {
       /*if (selectedTenantInfo &&
            selectedTenantInfo.providerIds && selectedTenantInfo.providerIds.length > 0) {
            const provider = new firebase.auth.OAuthProvider(selectedTenantInfo.providerIds[0]);
            // Set Additional Scopes
            provider.addScope('test_scope');

            // Set Query Parameters
            provider.setCustomParameters({
                "realm": "Underworld",
                login_hint: selectedTenantInfo.email || undefined,
            })
            auth.signInWithRedirect(provider);
        } else {
            reject("Couldn't resolve tenant from tenant info")
        }*/
        // Select identity Provider
        const provider = new firebase.auth.OAuthProvider('oidc.oauth-idp-config');
        auth.signInWithRedirect(provider);
    });
}
const ciapInstance = new ciap.Authentication(handler);
ciapInstance.start();