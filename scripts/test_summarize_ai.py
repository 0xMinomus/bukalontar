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

    def test_accepted_but_unreadable_is_not_consistent_reading(self):
        response = {'paleography': {'transliteration': 'Tidak terbaca.', 'aksara_bali_detected': True}, 'recipes': []}
        rows = [{'variant': 'check', 'model': 'check', 'image_sha256': 'a', 'elapsed_ms': 100, 'response': response} for _ in range(2)]
        report = summarize(rows, {'a': {'aksara_bali_detected': True}})[0]
        self.assertEqual(report['accepted_balinese_runs'], 2)
        self.assertEqual(report['detected_runs_with_no_readable_text'], 2)
        self.assertIsNone(report['transliteration_pair_agreement'])

    def test_provider_errors_are_not_correct_script_decisions(self):
        accepted = {'paleography': {'transliteration': 'bali', 'aksara_bali_detected': True}, 'recipes': []}
        rejected = {'paleography': {'transliteration': 'Tidak terdeteksi.', 'aksara_bali_detected': False}, 'recipes': []}
        rows = [
            {'variant': 'check', 'model': 'check', 'image_sha256': 'bali', 'elapsed_ms': 100, 'response': accepted},
            {'variant': 'check', 'model': 'check', 'image_sha256': 'bali', 'elapsed_ms': 100, 'error': 'timeout'},
            {'variant': 'check', 'model': 'check', 'image_sha256': 'other', 'elapsed_ms': 100, 'response': rejected},
            {'variant': 'check', 'model': 'check', 'image_sha256': 'other', 'elapsed_ms': 100, 'error': 'upstream'},
        ]
        report = summarize(rows, {'bali': {'aksara_bali_detected': True}, 'other': {'aksara_bali_detected': False}})[0]
        self.assertEqual((report['expected_balinese_runs'], report['accepted_balinese_runs']), (2, 1))
        self.assertEqual((report['expected_non_balinese_runs'], report['rejected_non_balinese_runs']), (2, 1))


if __name__ == '__main__':
    unittest.main()
