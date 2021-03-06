#include <thread>
#include "queries_api.h"

using std::thread;

TEST_F(QueriesApi, QueriesStandardTestSet) {
	FillDefaultNamespace(0, 2500, 20);
	FillDefaultNamespace(2500, 2500, 0);
	FillCompositeIndexesNamespace(0, 1000);
	FillTestSimpleNamespace();
	FillComparatorsNamespace();

	CheckStandartQueries();
	CheckAggregationQueries();
	CheckSqlQueries();
	CheckCompositeIndexesQueries();
	CheckComparatorsQueries();
	CheckDistinctQueries();

	int itemsCount = 0;
	InsertedItemsByPk& items = insertedItems[default_namespace];
	for (auto it = items.begin(); it != items.end();) {
		Error err = rt.reindexer->Delete(default_namespace, it->second);
		EXPECT_TRUE(err.ok()) << err.what();
		it = items.erase(it);
		if (++itemsCount == 4000) break;
	}

	FillDefaultNamespace(0, 500, 0);
	FillDefaultNamespace(0, 1000, 5);

	itemsCount = 0;
	// items = insertedItems[default_namespace];
	for (auto it = items.begin(); it != items.end();) {
		Error err = rt.reindexer->Delete(default_namespace, it->second);
		EXPECT_TRUE(err.ok()) << err.what();
		it = items.erase(it);
		if (++itemsCount == 5000) break;
	}

	for (size_t i = 0; i < 5000; ++i) {
		auto itToRemove = items.begin();
		if (itToRemove != items.end()) {
			Error err = rt.reindexer->Delete(default_namespace, itToRemove->second);
			EXPECT_TRUE(err.ok()) << err.what();
			items.erase(itToRemove);
		}
		FillDefaultNamespace(rand() % 100, 1, 0);

		if (!items.empty()) {
			itToRemove = items.begin();
			std::advance(itToRemove, rand() % std::min(100, int(items.size())));
			if (itToRemove != items.end()) {
				Error err = rt.reindexer->Delete(default_namespace, itToRemove->second);
				EXPECT_TRUE(err.ok()) << err.what();
				items.erase(itToRemove);
			}
		}
	}

	for (auto it = items.begin(); it != items.end();) {
		Error err = rt.reindexer->Delete(default_namespace, it->second);
		EXPECT_TRUE(err.ok()) << err.what();
		it = items.erase(it);
	}

	FillDefaultNamespace(3000, 1000, 20);
	FillDefaultNamespace(1000, 500, 00);
	FillCompositeIndexesNamespace(1000, 1000);
	FillComparatorsNamespace();

	CheckStandartQueries();
	CheckAggregationQueries();
	CheckSqlQueries();
	CheckCompositeIndexesQueries();
	CheckComparatorsQueries();
	CheckDistinctQueries();
}

TEST_F(QueriesApi, TransactionStress) {
	vector<thread> pool;
	FillDefaultNamespace(0, 350, 20);
	FillDefaultNamespace(3500, 350, 0);
	std::atomic_uint current_size;
	current_size = 350;
	uint32_t stepSize = 1000;

	for (size_t i = 0; i < 4; i++) {
		pool.push_back(thread([this, i, &current_size, stepSize]() {
			size_t start_pos = i * stepSize;
			if (i % 2 == 0) {
				uint32_t steps = 10;
				for (uint32_t j = 0; j < steps; ++j) {
					current_size += stepSize / steps;
					AddToDefaultNamespace(start_pos, start_pos + stepSize / steps, 20);
					start_pos = start_pos + stepSize / steps;
				}
			} else if (i % 2 == 1) {
				uint32_t oldsize = current_size.load();
				current_size += oldsize;
				FillDefaultNamespaceTransaction(current_size, start_pos + oldsize, 10);
			};
		}));
	}

	for (auto& tr : pool) {
		tr.join();
	}
}

TEST_F(QueriesApi, QueriesSqlGenerate) {
	const auto check = [](const string& sql) {
		Query q;
		q.FromSQL(sql);
		EXPECT_EQ(sql, q.GetSQL());
	};

	check("SELECT ID,Year,Genre FROM test_namespace WHERE year > '2016' ORDER BY year DESC LIMIT 10000000");

	check(
		"SELECT ID FROM test_namespace WHERE name LIKE 'something' AND (genre IN ('1','2','3') AND year > '2016') OR age IN "
		"('1','2','3','4') LIMIT 10000000");
}

std::vector<int> generateForcedSortOrder(int maxValue, size_t size) {
	std::set<int> res;
	while (res.size() < size) res.insert(rand() % maxValue);
	return {res.cbegin(), res.cend()};
}

TEST_F(QueriesApi, ForcedSortOffsetTest) {
	FillForcedSortNamespace();
	for (size_t i = 0; i < 100; ++i) {
		const auto forcedSortOrder =
			generateForcedSortOrder(forcedSortOffsetMaxValue * 1.1, rand() % static_cast<int>(forcedSortOffsetNsSize * 1.1));
		const size_t offset = rand() % static_cast<size_t>(forcedSortOffsetNsSize * 1.1);
		const size_t limit = rand() % static_cast<size_t>(forcedSortOffsetNsSize * 1.1);
		const bool desc = rand() % 2;
		// Single column sort
		auto expectedResults = ForcedSortOffsetTestExpectedResults(offset, limit, desc, forcedSortOrder, First);
		ExecuteAndVerify(forcedSortOffsetNs,
						 Query(forcedSortOffsetNs).Sort(kFieldNameColumnHash, desc, forcedSortOrder).Offset(offset).Limit(limit),
						 kFieldNameColumnHash, expectedResults);
		expectedResults = ForcedSortOffsetTestExpectedResults(offset, limit, desc, forcedSortOrder, Second);
		ExecuteAndVerify(forcedSortOffsetNs,
						 Query(forcedSortOffsetNs).Sort(kFieldNameColumnTree, desc, forcedSortOrder).Offset(offset).Limit(limit),
						 kFieldNameColumnTree, expectedResults);
		// Multicolumn sort
		const bool desc2 = rand() % 2;
		auto expectedResultsMult = ForcedSortOffsetTestExpectedResults(offset, limit, desc, desc2, forcedSortOrder, First);
		ExecuteAndVerify(forcedSortOffsetNs,
						 Query(forcedSortOffsetNs)
							 .Sort(kFieldNameColumnHash, desc, forcedSortOrder)
							 .Sort(kFieldNameColumnTree, desc2)
							 .Offset(offset)
							 .Limit(limit),
						 kFieldNameColumnHash, expectedResultsMult.first, kFieldNameColumnTree, expectedResultsMult.second);
		expectedResultsMult = ForcedSortOffsetTestExpectedResults(offset, limit, desc, desc2, forcedSortOrder, Second);
		ExecuteAndVerify(forcedSortOffsetNs,
						 Query(forcedSortOffsetNs)
							 .Sort(kFieldNameColumnTree, desc2, forcedSortOrder)
							 .Sort(kFieldNameColumnHash, desc)
							 .Offset(offset)
							 .Limit(limit),
						 kFieldNameColumnHash, expectedResultsMult.first, kFieldNameColumnTree, expectedResultsMult.second);
	}
}
