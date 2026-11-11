import unittest

from summarize_ai import edit_distance, summarize


class EvaluationMetricsTest(unittest.TestCase):
    def test_abstention_is_not_counted_as_consistent_reading(self):
        response = {
            'paleography': {'transliteration': '[Lembar 1]\n[tidak terbaca]', 'aksara_bali_detected': True},
            'recipes': [],
        }
        rows = [{'variant': 'check', 'model': 'check', 'image_sha256': 'a', 'elapsed_ms': 100, 'response': response} for _ in range(3)]
        report = summarize(rows, {})[0]
        self.assertEqual(report['detected_runs_with_no_readable_text'], 3)
        self.assertIsNone(report['transliteration_pair_agreement'])
        self.assertIsNone(report['plant_set_pair_jaccard'])
        self.assertIsNone(report['mean_character_error_rate'])

    def test_character_error_rate_requires_a_reference(self):
        self.assertEqual(edit_distance('bali', 'balu'), 1)
        self.assertEqual(edit_distance('', 'bali'), 4)
        row = {'variant': 'check', 'model': 'check', 'image_sha256': 'a', 'elapsed_ms': 100,
               'response': {'paleography': {'transliteration': 'balu', 'aksara_bali_detected': True}, 'recipes': []}}
        report = summarize([row], {'a': {'transliteration': 'bali'}})[0]
        self.assertEqual(report['mean_character_error_rate'], 0.25)
        self.assertEqual(report['transcription_references_scored'], 1)


if __name__ == '__main__':
    unittest.main()
