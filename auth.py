import requests
import re
import json
from html import unescape
import os

email = os.environ['EMAIL']
password = os.environ['PASS']

with requests.Session() as s:
    r = s.get('https://intra.epitech.eu')

    m = re.search('https:\/\/login.microsoftonline([^"]+)', r.text)
    url = 'https://login.microsoftonline' + m.group(1)

    r = s.get(url)
    m = re.search("Constants.CONTEXT = '([^']+)", r.text)
    urlRealm = 'https://login.microsoftonline.com/common/userrealm/?user=' + email + '&api-version=2.1&stsRequest=' + m.group(1) + '&checkForMicrosoftAccount=true';
    realmDatas = json.loads(requests.get(urlRealm).text)
    realmUrl = realmDatas['AuthURL']
    epitechLoginPage = s.get(realmUrl).text
    m = re.search('action="/adfs([^"]+)', epitechLoginPage)
    loginUrl = 'https://sts.epitech.eu/adfs' + m.group(1)
    r = s.post(loginUrl, data={'UserName': email, 'Password': password, 'Kmsi': 'true', 'AuthMethod': 'FormsAuthentication'})
    m = re.search('action="([^"]+)', r.text)
    microsoftLoginUrl = m.group(1)

    m = re.search('name="wa" value="([^"]+)', r.text)
    wa = m.group(1)
    m = re.search('name="wresult" value="([^"]+)', r.text)
    wresult = m.group(1)
    m = re.search('name="wctx" value="([^"]+)', r.text)
    wctx = m.group(1)

    microsoftLoginResult = s.post(microsoftLoginUrl, data={'wa': wa, 'wresult': unescape(wresult), 'wctx': unescape(wctx)}, allow_redirects=True)
    body = microsoftLoginResult.text

    if body.find("consent_accept_form") != -1:
        print ("Consent needed. Granting...\n")

    cookies = s.cookies.get_dict()

    print(s.get("https://intra.epitech.eu/?format=json").text)
