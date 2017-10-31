# built-in libs
import json
import random
import sys
from multiprocessing import Pool

# side libs
import bs4
import requests


def parse_site_urls(url: str) -> list:
    host = url[:url.find('/', 10)]
    r = requests.get(url)
    soup = bs4.BeautifulSoup(r.text, 'html.parser')
    raw_urls = [tag.attrs['href'] for tag in soup.find_all('a', attrs={'class': 'hotel_name_link url'})]
    return [host + url.strip('\n').split('\n')[0] for url in raw_urls]


def get_english_title(url: str) -> str:
    r = requests.get(url)
    soup = bs4.BeautifulSoup(r.text, 'html.parser')
    return str(soup.find('h2', attrs={'id': 'hp_hotel_name'}).string).strip('\n')


def get_hotel_data(hotel_id: int, url: str) -> dict:
    r = requests.get(url)
    soup = bs4.BeautifulSoup(r.text, 'html.parser')
    title_ru = str(soup.find('h2', attrs={'id': 'hp_hotel_name'}).string).strip('\n')

    url_eng = soup.find('link', attrs={'hreflang': 'en-us'}).attrs['href']
    title_en = get_english_title(url_eng)

    desc_tag = soup.find('div', attrs={'id': 'summary'}).contents
    description = '\n'.join(map(lambda x: x.text, filter(lambda x: isinstance(x, bs4.Tag) and x.name == 'p', desc_tag)))

    loc = soup.find('span', attrs={'class': 'hp_address_subtitle jq_tooltip'})
    geo = list(map(float, loc.attrs['data-bbox'].split(",")))
    location = {'lat': (geo[1] + geo[3]) / 2, 'lon': (geo[0] + geo[2]) / 2}
    address = strip_until_done(''.join(loc.contents))

    rooms = []
    available_rooms = []

    rooms_list = [str(tag.contents[2]).strip('\n') for tag in soup.find_all('a', attrs={'class': 'jqrt'})]
    rooms_extra = [tag.text.split(':')[1] for tag in
                   soup.find_all('p', attrs={'class': 'hp_rt_lightbox_facilities'})]

    m = min(len(rooms_list), len(rooms_extra))
    rooms_list = rooms_list[:m]
    rooms_extra = rooms_extra[:m]

    for idx, r in enumerate(rooms_list):
        beds = random.randint(1, 10)
        capacity = beds + random.randint(0, 5)
        rooms.append({
            'id': idx,
            'title': r,
            'type': 'room',
            'beds': beds,
            'capacity': capacity,
            'facilities': list(map(strip_until_done, rooms_extra[idx].split(',')))
        })

        available_rooms.append({
            'room_id': idx,
            'count': random.randint(0, 5)
        })

    amenities = [x.contents[-1].strip('\n') for x in soup.find_all('h5')]
    hotel = {
        'id': hotel_id,
        'title_ru': title_ru,
        'title_en': title_en,
        'description': description,
        'rooms': rooms,
        'location': location,
        'address': address,
        'amenities': amenities,
        'available_rooms': available_rooms
    }

    print(url + ' ok')
    return hotel


def get_hotel_data_wrapper(x):
    try:
        return get_hotel_data(x[0], x[1])
    except Exception as e:
        print(x[1] + ' fail')
        return {x[0]: 'caught error!: %s' % e}


def strip_until_done(s):
    ret = s
    while True:
        old = ret
        ret = ret.strip(' ')
        ret = ret.strip('\n')
        ret = ret.strip('\r')
        if old == ret:
            break

    return ret


if __name__ == '__main__':
    hotel_urls = parse_site_urls(sys.argv[1])
    hotels = []

    pool = Pool(5)
    results = pool.map(get_hotel_data_wrapper, enumerate(hotel_urls))
    with open('data.json', 'w') as f:
        f.write(json.dumps(results))
