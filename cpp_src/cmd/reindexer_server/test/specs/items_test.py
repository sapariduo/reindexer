import random
from specs import BaseTest


class ItemsTest(BaseTest):
    def setUp(self):
        super().setUp()

        self.current_db = self.test_db
        status, body = self.api_create_db(self.current_db)

        self.assertEqual(True, status == 200, body)

        self.current_ns = self.test_ns
        status, body = self.api_create_namespace(
            self.current_db, self.current_ns)

        self.assertEqual(True, status == 200, body)

    def test_get_items(self):
        """Should be able to get items list"""

        status, body = self.api_get_items(self.current_db, self.current_ns)
        self.assertEqual(True, status == 200, body)

    def test_create_item(self):
        """Should be able to create item"""

        index_count = 2
        index_array_of_dicts = self.helper_index_array_construct(index_count)

        for i in range(0, index_count):
            status, body = self.api_create_index(
                self.current_db, self.test_ns, index_array_of_dicts[i])
            self.assertEqual(True, status == 200, body)

        item_body = self.helper_item_construct()

        status, body = self.api_create_item(
            self.current_db, self.current_ns, item_body)
        self.assertEqual(True, status == 200, body)

        status, body = self.api_get_items(self.current_db, self.current_ns)
        self.assertEqual(True, status == 200, body)
        self.assertEqual(True, item_body in body['items'], body)

    def test_create_huge_item(self):
        """Should be able to create huge items"""

        index_count = 2
        index_array_of_dicts = self.helper_index_array_construct(index_count)

        for i in range(0, index_count):
            status, body = self.api_create_index(
                self.current_db, self.test_ns, index_array_of_dicts[i])
            self.assertEqual(True, status == 200, body)

        for i in range(0, 10):
            item_body = self.helper_item_construct(5, i, i * 10000)

            status, body = self.api_create_item(
                self.current_db, self.current_ns, item_body)
            self.assertEqual(True, status == 200, body)

            status, body = self.api_get_items(self.current_db, self.current_ns)
            self.assertEqual(True, status == 200, body)
            self.assertEqual(True, item_body in body['items'], body)

    def test_update_item(self):
        """Should be able to update item"""

        count = 5
        index_array_of_dicts = self.helper_index_array_construct(count)

        for i in range(0, count):
            status, body = self.api_create_index(
                self.current_db, self.test_ns, index_array_of_dicts[i])
            self.assertEqual(True, status == 200, body)

        item_body = self.helper_item_construct(count)

        status, body = self.api_create_item(
            self.current_db, self.current_ns, item_body)
        self.assertEqual(True, status == 200, body)

        status, body = self.api_get_items(self.current_db, self.current_ns)
        self.assertEqual(True, status == 200, body)
        self.assertEqual(True, item_body in body['items'], body)

        item_body_last_key = 'test_' + str(count)
        new_item_body = item_body
        rand_int = random.randint(0xFF, 0x7FFFFFFF)
        new_item_body[item_body_last_key] = rand_int

        status, body = self.api_update_item(
            self.current_db, self.current_ns, new_item_body)
        self.assertEqual(True, status == 200, body)

        status, body = self.api_get_sorted_items(
            self.current_db, self.current_ns, item_body_last_key, 'desc')
        self.assertEqual(True, status == 200, body)
        self.assertEqual(True, new_item_body in body['items'], new_item_body)

    def test_delete_item(self):
        """Should be able to delete an item"""

        count = 5
        index_array_of_dicts = self.helper_index_array_construct(count)

        for i in range(0, count):
            status, body = self.api_create_index(
                self.current_db, self.test_ns, index_array_of_dicts[i])
            self.assertEqual(True, status == 200, body)

        item_body = self.helper_item_construct(count)

        status, body = self.api_create_item(
            self.current_db, self.current_ns, item_body)
        self.assertEqual(True, status == 200, body)

        status, body = self.api_get_items(self.current_db, self.current_ns)
        self.assertEqual(True, status == 200, body)
        self.assertEqual(True, item_body in body['items'], body)

        status, body = self.api_delete_item(
            self.current_db, self.current_ns, item_body)
        self.assertEqual(True, status == 200, body)

        status, body = self.api_get_items(self.current_db, self.current_ns)
        self.assertEqual(True, status == 200, body)
        self.assertEqual(True, item_body not in body['items'], body)

    def test_get_items_pagination(self):
        """Should be able to get item list with limit and offset parameters"""

        index_count = 2
        index_array_of_dicts = self.helper_index_array_construct(index_count)

        for i in range(0, index_count):
            status, body = self.api_create_index(
                self.current_db, self.test_ns, index_array_of_dicts[i])
            self.assertEqual(True, status == 200, body)

        items_count = 20
        items = self.helper_item_array_construct(items_count)

        for item_body in items:
            status, body = self.api_create_item(
                self.current_db, self.current_ns, item_body)
            self.assertEqual(True, status == 200, body)

        limit = 2
        offset = 10
        item_index = offset
        status, body = self.api_get_paginated_items(
            self.current_db, self.current_ns, limit, offset)
        self.assertEqual(True, status == 200, body)
        self.assertEqual(True, limit == len(body['items']), body)
        self.assertEqual(True, items[item_index] in body['items'], body)
        self.assertEqual(True, body['total_items'] == items_count, body)

    def test_get_items_sort_asc_sort_dir_param(self):
        """Should be able to get asc-sorted item list using sort_field and sort_order=asc parameters"""

        index_count = 2
        index_array_of_dicts = self.helper_index_array_construct(index_count)

        for i in range(0, index_count):
            status, body = self.api_create_index(
                self.current_db, self.test_ns, index_array_of_dicts[i])
            self.assertEqual(True, status == 200, body)

        items_count = 20
        items = self.helper_item_array_construct(items_count)

        for item_body in items:
            status, body = self.api_create_item(
                self.current_db, self.current_ns, item_body)
            self.assertEqual(True, status == 200, body)

        sort_field = self.helper_items_first_key_of_item(items)
        sort_dir = "asc"

        status, body = self.api_get_sorted_items(
            self.current_db, self.current_ns, sort_field, sort_dir)
        self.assertEqual(True, status == 200, body)

        self.assertEqual(True, body['items'][0][sort_field]
                         < body['items'][-1][sort_field], body)

    def test_get_items_sort_desc_sort_dir_param(self):
        """Should be able to get desc-sorted item list using sort_field and sort_order=desc parameters"""

        index_count = 2
        index_array_of_dicts = self.helper_index_array_construct(index_count)

        for i in range(0, index_count):
            status, body = self.api_create_index(
                self.current_db, self.test_ns, index_array_of_dicts[i])
            self.assertEqual(True, status == 200, body)

        items_count = 20
        items = self.helper_item_array_construct(items_count)

        for item_body in items:
            status, body = self.api_create_item(
                self.current_db, self.current_ns, item_body)
            self.assertEqual(True, status == 200, body)

        sort_field = self.helper_items_first_key_of_item(items)
        sort_dir = "desc"

        status, body = self.api_get_sorted_items(
            self.current_db, self.current_ns, sort_field, sort_dir)
        self.assertEqual(True, status == 200, body)

        self.assertEqual(True, body['items'][0][sort_field]
                         > body['items'][-1][sort_field], body)

    def test_get_items_sort_wrong_sort_any_or_empty_param(self):
        """Should be able to get asc-sorted item list using sort_field and any or empty sort_order parameters"""

        index_count = 2
        index_array_of_dicts = self.helper_index_array_construct(index_count)

        for i in range(0, index_count):
            status, body = self.api_create_index(
                self.current_db, self.test_ns, index_array_of_dicts[i])
            self.assertEqual(True, status == 200, body)

        items_count = 20
        items = self.helper_item_array_construct(items_count)

        for item_body in items:
            status, body = self.api_create_item(
                self.current_db, self.current_ns, item_body)
            self.assertEqual(True, status == 200, body)

        sort_field = self.helper_items_first_key_of_item(items)
        sort_dir = 'wrong'

        status, body = self.api_get_sorted_items(
            self.current_db, self.current_ns, sort_field, sort_dir)
        self.assertEqual(True, status == 200, body)

        self.assertEqual(True, body['items'][0][sort_field]
                         < body['items'][-1][sort_field], body)

        sort_dir = ''

        status, body = self.api_get_sorted_items(
            self.current_db, self.current_ns, sort_field, sort_dir)
        self.assertEqual(True, status == 200, body)

        self.assertEqual(True, body['items'][0][sort_field]
                         < body['items'][-1][sort_field], body)